package rpc

import (
	"github.com/aaronzjc/mu/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCommanderClient(addr string) (pb.CommanderClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewCommanderClient(conn), nil
}

func NewAgentClient(addr string) (pb.AgentClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewAgentClient(conn), nil
}
