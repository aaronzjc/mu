package commander

import (
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"mu/internal/svc/rpc"
	"mu/internal/svc/schedule"
	"mu/internal/util/logger"
	"net"
)

type CommanderServer struct{}

func (commander *CommanderServer) UpdateCron(ctx context.Context, req *rpc.Cron) (*rpc.CronRes, error) {
	schedule.JobSchedule.UpdateJob(req.Site)
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
