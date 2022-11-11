package commander

import (
	"context"

	"github.com/aaronzjc/mu/internal/constant"
	"github.com/aaronzjc/mu/internal/infra/cache"
	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/pkg/logger"
)

type CommanderServer struct {
	pb.UnimplementedCommanderServer
}

var _ pb.CommanderServer = &CommanderServer{}

func (commander *CommanderServer) UpdateCron(ctx context.Context, req *pb.Cron) (*pb.CronRes, error) {
	redis := cache.Get()
	redis.LPush(constant.JobVisor, req.Site)
	logger.Info("Rpc UpdateCron [site = %s] success !", req.Site)
	return &pb.CronRes{Success: true}, nil
}

func NewCommanderServer() *CommanderServer {
	return &CommanderServer{}
}
