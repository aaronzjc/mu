package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/commander"
	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/aaronzjc/mu/internal/infra/cache"
	"github.com/aaronzjc/mu/internal/infra/db"
	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/pkg/logger"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	Identifier string // 集群中的唯一标识
	IsLeader   bool   // 是否是主节点
)

func SetupCommander(ctx *cli.Context) error {
	var err error
	if ctx.String("config") == "" {
		return errors.New("invalid config option, use -h get full doc")
	}
	// 初始化项目配置
	configFile := ctx.String("config")
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		return errors.New("config file not found")
	}
	var conf *config.Config
	if conf, err = config.LoadConfig(configFile); err != nil {
		return err
	}
	// 初始化日志组件
	if err = logger.Setup("commander", "/var/log/mu-commander.log"); err != nil {
		return err
	}
	// 初始化DB
	if err = db.Setup(conf, &gorm.Config{}); err != nil {
		return err
	}
	// 初始化缓存
	if err = cache.Setup(&conf.Redis); err != nil {
		return err
	}
	// 初始化任务调度
	commander.SetupSchedule()
	// 设置调试模式
	if conf.Env != "prod" {
		logger.SetLevel(conf.Log.Level)
	}

	return nil
}

func RunCommander(ctx *cli.Context) error {
	conf := config.Get()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Commander.Port))
	if err != nil {
		logger.Fatal("bind socket failed")
	}

	var opts []grpc.ServerOption
	rpcServer := grpc.NewServer(opts...)
	rpcServer.RegisterService(&pb.Commander_ServiceDesc, commander.NewCommanderServer())

	go rpcServer.Serve(listener)
	logger.Info("[START] commander listen at :", conf.Commander.Port)

	// 进行选举管理
	ctxx, cancelManager := context.WithCancel(context.Background())
	go commanderManager(ctxx)

	// 监听关闭信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM)
	<-sig

	// 关闭服务
	cancelManager() // 通知选举等停止
	rpcServer.GracefulStop()
	logger.Info("[STOP] commander stop done")
	return nil
}

// 公共职责：监听master状态，如果master挂了，重新选一个master。
// 0. 获取一个节点标识
// 1. 如果master正常，则不做什么操作。
// 2. 如果master不正常，没有更新自己的健康状态
// 2.1. 如果当前节点是master，那么注销之前注册的定时任务，和队列监听
// 2.2. 如果当前节点不是master，那么设置自己为新的master，然后注册任务
func commanderManager(ctx context.Context) {
	redis := cache.Get()
	// 初始化当前节点ID
	idRes := redis.Incr(constant.IdMachine)
	Identifier = strconv.Itoa(int(idRes.Val()))

	t := time.NewTicker(time.Second * 3)
	defer t.Stop()

	subCtx, cancel := context.WithCancel(ctx)
	for {
		select {
		case <-ctx.Done():
			logger.Info("receive ctx cancel, gracefully stop manager")
			cancel() // cancel child ctx
			return
		case <-t.C:
			masterId := redis.Get(constant.Election).Val()

			var needElection bool
			if masterId == "" {
				needElection = true
				logger.Info("need re-election")
			}

			if needElection {
				if ok := redis.SetNX(constant.Election, Identifier, time.Second*10); ok.Val() {
					logger.Info("election done, current master is " + Identifier)
					masterId = Identifier
					// 只有之前不是leader，新当上leader才走这里
					if !IsLeader {
						IsLeader = true
						// 选举成功，恭喜当上老大
						go masterDuty(subCtx)
					}
				} else {
					logger.Info("election failed")
				}
			}

			if IsLeader && masterId != Identifier {
				// 不好意思，你不再担任老大了
				IsLeader = false
				for i := 0; i < 5; i++ {
					if err := commander.JobSchedule.TruncateJobs(); err == nil {
						logger.Info("Done truncate current jobs")
						break
					}
				}
			}
		}
	}
}

// Master的职责
// 1. 注册定时任务
// 2. 处理队列的任务，更新定时任务。
// 3. 上报自己的健康状态
func masterDuty(ctx context.Context) {
	// 注册定时任务
	commander.JobSchedule.InitJobs()

	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if !IsLeader {
				// 再次判断是否是master，可能中途出现状况，执行了主从切换
				logger.Info("i am not master, id = ", Identifier)
				break
			}

			redis := cache.Get()

			// 上报自己的状态
			redis.Set(constant.Election, Identifier, time.Second*10)
			logger.Debug(fmt.Sprintf("Health Set success , time at %v", time.Now()))

			for {
				data := redis.LPop(constant.JobVisor)
				if data.Val() == "" {
					logger.Debug("empty queue, break")
					break
				}
				logger.Info(fmt.Sprintf("receive update [site = %s] from queue", data.Val()))
				svc := service.NewSiteService(store.NewSiteRepo(), nil)
				sites, _ := svc.Get(ctx, &dto.Query{
					Query: "`key` = ?",
					Args:  []interface{}{data.Val()},
				})
				if len(sites) != 0 {
					crawJob := commander.NewCrawJob(sites[0])
					commander.JobSchedule.UpdateJob(sites[0].Key, sites[0].Cron, crawJob)
				} else {
					logger.Error("update site = " + data.Val() + " craw job err")
				}
			}
		}
	}
}
