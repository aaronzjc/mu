package crawler

import (
	"context"
	"crawler/internal/svc/lib"
	"crawler/internal/svc/rpc"
	"crawler/internal/util/logger"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

type AgentServer struct {}

func (agent *AgentServer) Craw(ctx context.Context, msg *rpc.Job) (*rpc.Result, error) {
	var wg sync.WaitGroup
	pageMap := make(map[string]lib.Page)
	s := lib.FSite(msg.Name)
	links, _ := s.BuildUrl()
	for _, link := range links {
		wg.Add(1)
		go func(link lib.Link) {
			page, err := s.CrawPage(link)
			if err != nil {
				logger.Error("craw page error, err %v .", err)
				return
			}
			logger.Info("craw page done %s .", link.Url)
			pageMap[link.Tag] = page
			wg.Done()
		}(link)
	}

	wg.Wait()

	result := new(rpc.Result)
	result.T = time.Now().Format("2006-01-02 15:04:05")
	m := make(map[string]*rpc.Result_HotList)
	for tag, p := range pageMap {
		hotList := new(rpc.Result_HotList)
		var items []*rpc.Result_HotList_Item
		for _, item := range p.List {
			items = append(items, &rpc.Result_HotList_Item{
				Title: item.Title,
				Rank:  float32(item.Rank),
				Url:   item.OriginUrl,
			})
		}
		hotList.Item = items
		m[tag] = hotList
	}
	result.Map = m

	return result, nil
}

func (agent *AgentServer) Check(ctx context.Context, msg *rpc.Ping) (*rpc.Pong, error) {
	logger.Info("receive check")
	return &rpc.Pong{
		Pong: msg.Ping,
	}, nil
}

func RegisterRpcServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("bind socket failed")
	}

	var opts []grpc.ServerOption
	rpcServer := grpc.NewServer(opts...)
	rpc.RegisterAgentServer(rpcServer, &AgentServer{})
	logger.Info("server is listening on :7990")
	logger.Fatal(rpcServer.Serve(lis))
}