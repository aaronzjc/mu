package agent

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/aaronzjc/mu/internal/core/site"
	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/pkg/helper"
	"github.com/aaronzjc/mu/pkg/logger"
)

type AgentServer struct {
	pb.UnimplementedAgentServer
}

var _ pb.AgentServer = &AgentServer{}

func (agent *AgentServer) Craw(ctx context.Context, msg *pb.Job) (*pb.Result, error) {
	var wg sync.WaitGroup

	pageMap := make(map[string]site.Page)
	headers := make(map[string]string)

	h := msg.Headers
	for _, v := range h {
		headers[v.Key] = v.Val
	}
	spider, ok := site.SiteMap[msg.Name]
	if !ok {
		return nil, errors.New("not supported site " + msg.Name)
	}

	links, _ := spider.BuildUrl()
	for _, link := range links {
		wg.Add(1)
		go func(link site.Link) {
			defer wg.Done()
			page, err := spider.CrawPage(link, headers)
			if err != nil {
				logger.Error("craw page error, err " + err.Error())
				return
			}
			logger.Info("craw page done, link = " + link.Url)
			pageMap[link.Tag] = page
		}(link)
	}

	wg.Wait()

	result := new(pb.Result)
	result.T = helper.CurrentTimeStr()
	result.HotMap = make(map[string]string)
	for tag, page := range pageMap {
		res, _ := json.Marshal(page.List)
		result.HotMap[tag] = string(res)
	}
	return result, nil
}

func (agent *AgentServer) Check(ctx context.Context, msg *pb.Ping) (*pb.Pong, error) {
	logger.Info("receive health check")
	return &pb.Pong{
		Pong: msg.Ping,
	}, nil
}

func NewAgentServer() *AgentServer {
	return &AgentServer{}
}
