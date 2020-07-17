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
	"time"
)

const idMachine = "is_machine"
const jobVisor = "job_visor"
const election = "election"

var (
	id	int
	isLeader bool
)

type CommanderServer struct{}

func (commander *CommanderServer) UpdateCron(ctx context.Context, req *rpc.Cron) (*rpc.CronRes, error) {
	redis := cache.RedisConn()
	redis.LPush(jobVisor, req.Site)
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
func InitCommander() bool {
	// 选择一个主节点
	redis := cache.RedisConn()

	// 获取当前节点ID
	idRes := redis.Incr(idMachine)
	id = int(idRes.Val())

	// 选择一个节点作为leader
	if ok := redis.SetNX(election, id, time.Minute); ok.Val() {
		isLeader = true
	}

	// 注册定时任务
	if isLeader {
		schedule.JobSchedule.InitJobs()
	}

	return isLeader
}

// 监听队列里面任务
func JobVisor() {
	if !isLeader {
		return
	}

	redis := cache.RedisConn()
	for {
		data := redis.LPop(jobVisor)
		if data.Val() == "" {
			time.Sleep(time.Second * 5)
			continue
		}
		schedule.JobSchedule.UpdateJob(data.Val())
		logger.Info("Rpc UpdateCron [site = %s] success !", data)
	}
}

// 管理master的健康状态
func MasterVisor() {
	// TODO: 管理
}