package commander

import (
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"mu/internal/svc/rpc"
	"mu/internal/svc/schedule"
	"mu/internal/util/cache"
	"mu/internal/util/logger"
	"net"
	"strconv"
	"time"
)

const idMachine = "id_machine"
const jobVisor = "job_visor"
const election = "election"

var (
	id	string
	isLeader bool
)

type CommanderServer struct{}

func (commander *CommanderServer) UpdateCron(ctx context.Context, req *rpc.Cron) (*rpc.CronRes, error) {
	redis := cache.RedisConn()
	redis.LPush(jobVisor, req.Site)
	logger.Info("Rpc UpdateCron [site = %s] success !", req.Site)
	return &rpc.CronRes{Success: true}, nil
}

func (commander *CommanderServer) Debug(ctx context.Context, req *rpc.Empty) (*rpc.DebugRes, error) {
	data := schedule.Debug()
	js, _ := json.Marshal(data)
	logger.Info("Rpc Debug success !")
	return &rpc.DebugRes{Res: string(js)}, nil
}

func RegisterRpcServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("bind socket failed")
	}

	var opts []grpc.ServerOption
	rpcServer := grpc.NewServer(opts...)
	rpc.RegisterCommanderServer(rpcServer, &CommanderServer{})
	logger.Info("Commanderserver is listening on :7970")
	logger.Fatal(rpcServer.Serve(lis))
}

// 初始化，分配一个ID，并且选择一个主节点
func InitCommander() {
	redis := cache.RedisConn()

	// 初始化当前节点ID
	idRes := redis.Incr(idMachine)
	id = strconv.Itoa(int(idRes.Val()))
	logger.Info("current node id is %s", id)

	// 检测&选举&初始化
	go ManageMaster()
}

// Master的职责
// 1. 注册定时任务
// 2. 处理队列的任务，更新定时任务。
// 3. 上报自己的健康状态
func MasterDuty() {
	// 注册定时任务
	schedule.JobSchedule.InitJobs()

	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		<- t.C
		if !isLeader {
			// 再次判断是否是master，可能中途出现状况，执行了主从切换
			logger.Info("%s no more master, done own duty", id)
			break
		}

		redis := cache.RedisConn()

		// 上报自己的状态
		redis.Set(election, id, time.Second * 10)
		logger.Info("Health Set success , time at %v", time.Now())

		for {
			data := redis.LPop(jobVisor)
			if data.Val() == "" {
				logger.Info("empty queue, break")
				break
			}
			schedule.JobSchedule.UpdateJob(data.Val())
			logger.Info("Rpc UpdateCron [site = %s] success !", data)
		}
	}
}

// 公共职责：监听master状态，如果master挂了，重新选一个master。
// 1. 如果master正常，则不做什么操作。
// 2. 如果master不正常，没有更新自己的健康状态
// 2.1. 如果当前节点是master，那么注销之前注册的定时任务，和队列监听
// 2.2. 如果当前节点不是master，那么设置自己为新的master，然后注册任务
func ManageMaster() {
	t := time.NewTicker(time.Second * 3)
	defer t.Stop()
	for {
		<- t.C
		redis := cache.RedisConn()
		masterId := redis.Get(election).Val()

		var needElection bool
		if masterId == "" {
			needElection = true
			logger.Info("need re-election")
		}

		if needElection {
			if ok := redis.SetNX(election, id, time.Second * 10); ok.Val() {
				logger.Info("election done, current master is %s", id)
				masterId = id
				isLeader = true
				// 选举成功，恭喜当上老大
				go MasterDuty()
			}
		}

		if isLeader && masterId != id {
			// 不好意思，你不再担任老大了
			isLeader = false
			for i := 0; i < 5; i++ {
				err := schedule.JobSchedule.TruncateJobs()
				if err == nil {
					logger.Info("Done truncate current jobs")
					break
				}
			}
		}
	}
}