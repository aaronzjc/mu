package internal

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaronzjc/mu/internal/agent"
	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/pkg/logger"
	"github.com/aaronzjc/mu/test"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func SetupAgent(ctx *cli.Context) error {
	err := logger.Setup("agent", "")
	if err != nil {
		return err
	}

	// 测试环境配置代理
	if os.Getenv("APP_ENV") != "prod" {
		test.SetupProxy()
	}
	return nil
}

func RunAgent(ctx *cli.Context) error {
	addr := ":7990"
	listener, err := net.Listen("tcp", addr) // no need to use config file
	if err != nil {
		logger.Fatal("bind socket failed")
	}

	var opts []grpc.ServerOption
	rpcServer := grpc.NewServer(opts...)
	rpcServer.RegisterService(&pb.Agent_ServiceDesc, agent.NewAgentServer())

	go rpcServer.Serve(listener)
	logger.Info("[START] agent listen at ", addr)

	// 监听关闭信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM)
	<-sig

	// 关闭服务
	rpcServer.GracefulStop()
	logger.Info("[STOP] agent stop done")
	return nil
}
