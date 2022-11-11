package rpc

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/pkg/logger"
	"google.golang.org/grpc"
)

type RpcClient struct {
	Conn   *grpc.ClientConn
	Client *pb.AgentClient
}

type RpcPool struct {
	Clients map[string]*RpcClient
	Lock    sync.RWMutex
}

func (r *RpcPool) Get(addr string) (*RpcClient, error) {
	r.Lock.RLock()
	rc, ok := r.Clients[addr]
	r.Lock.RUnlock()
	if ok {
		return rc, nil
	}

	client, err := r.Set(addr)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *RpcPool) Set(addr string) (*RpcClient, error) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	rc, ok := r.Clients[addr]
	if ok {
		return rc, nil
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logger.Error("connect error " + err.Error())
		return nil, errors.New("dial server " + addr + " failed")
	}

	client := pb.NewAgentClient(conn)
	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	r.Clients[addr] = &RpcClient{
		Conn:   conn,
		Client: &client,
	}

	return r.Clients[addr], nil
}

func (r *RpcPool) Release(addr string) bool {
	r.Lock.Lock()
	rc, ok := r.Clients[addr]
	r.Lock.Unlock()
	if !ok {
		return true
	}

	delete(r.Clients, addr)
	_ = rc.Conn.Close()

	return true
}
