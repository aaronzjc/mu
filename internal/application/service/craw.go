package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/core/rpc"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/pkg/logger"
)

type CrawService interface {
	PickAgent(context.Context, *dto.Site) (*dto.Node, error)
	Craw(context.Context, *dto.Site) error
}

type CrawServiceImpl struct {
	siteRepo repo.SiteRepo
	nodeRepo repo.NodeRepo
}

var _ CrawService = &CrawServiceImpl{}

func (c *CrawServiceImpl) PickAgent(ctx context.Context, site *dto.Site) (*dto.Node, error) {
	rand.Seed(time.Now().UnixNano())
	var nodes []model.Node

	q := &dto.Query{}
	if site.NodeOption == model.ByType {
		q.Query = "`type` = ? AND `ping` = ?"
		q.Args = []interface{}{model.ByType, model.PingOk}
	} else {
		if len(site.NodeHosts) == 0 {
			return nil, errors.New("no nodes configured")
		}
		q.Query = "`id` IN (?) AND `enable` = ? AND `ping` = ?"
		q.Args = []interface{}{site.NodeHosts, model.Enable, model.PingOk}
	}
	nodes, err := c.nodeRepo.Get(ctx, q)
	if err != nil {
		return nil, errors.New("get nodes failed")
	}
	if len(nodes) == 0 {
		return nil, errors.New("no nodes avaiable")
	}
	chosen := nodes[rand.Int()%len(nodes)]
	return (&dto.Node{}).FillByModel(chosen), nil
}

func (c *CrawServiceImpl) Craw(ctx context.Context, site *dto.Site) error {
	picked, err := c.PickAgent(ctx, site)
	if err != nil {
		return errors.New("no agent avaiable")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rpcClient, _ := rpc.NewAgentClient(picked.Addr)

	// read craw config
	var headers []*pb.Job_Header
	for _, v := range site.ReqHeaders {
		if v.Key == "" || v.Val == "" {
			continue
		}
		headers = append(headers, &pb.Job_Header{
			Key: v.Key,
			Val: v.Val,
		})
	}

	// do craw
	var result *pb.Result
	if result, err = rpcClient.Craw(ctx, &pb.Job{
		Name:    site.Key,
		Headers: headers,
	}); err != nil {
		logger.Error("remote craw err %v", err)
		return err
	}
	logger.Info("remote craw [" + site.Key + "] done")

	// save to cache
	var news dto.News
	news.T = result.T
	for tag, hotStr := range result.HotMap {
		json.Unmarshal([]byte(hotStr), &news.List)
		data, _ := json.Marshal(news)
		err := c.siteRepo.SaveNews(ctx, site.Key, tag, string(data))
		if err != nil {
			logger.Error("save new err " + err.Error())
		}
	}
	return nil
}

func NewCrawService() *CrawServiceImpl {
	return &CrawServiceImpl{
		siteRepo: store.NewSiteRepo(),
		nodeRepo: store.NewNodeRepo(),
	}
}
