package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/core/rpc"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/pkg/logger"
)

type NodeService interface {
	CheckNodes(context.Context, *rpc.RpcPool) error
	Upsert(context.Context, *dto.Node) error
	Del(context.Context, int) error
	Get(context.Context, *dto.Query) ([]*dto.Node, error)
}

type NodeServiceImpl struct {
	repo repo.NodeRepo
}

var _ NodeService = &NodeServiceImpl{}

func (s *NodeServiceImpl) CheckNodes(ctx context.Context, clientPool *rpc.RpcPool) error {
	nodes, err := s.repo.Get(ctx, &dto.Query{
		Query: "`enable` = ?",
		Args:  []interface{}{model.Enable},
	})
	if err != nil {
		return errors.New("get nodes err")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	for _, node := range nodes {
		var err error
		var client *rpc.RpcClient
		if client, err = clientPool.Get(node.Addr); err != nil {
			if node.Ping != model.PingFail {
				s.repo.Update(ctx, node, map[string]interface{}{
					"ping": model.PingFail,
				})
			}
			continue
		}
		ping := &pb.Ping{Ping: "ping"}
		if result, err := (*client.Client).Check(ctx, ping); err != nil || result.Pong != ping.Ping {
			logger.Error(fmt.Sprintf("rpc health check : [%s] ping error, err %v.", node.Name, err))
			if node.Ping != model.PingFail {
				s.repo.Update(ctx, node, map[string]interface{}{
					"ping": model.PingFail,
				})
			}
			continue
		}
		logger.Info(fmt.Sprintf("rpc health check : [%s] is online.", node.Name))
		if node.Ping != model.PingOk {
			s.repo.Update(ctx, node, map[string]interface{}{
				"ping": model.PingOk,
			})
		}
	}
	return nil
}

func (s *NodeServiceImpl) Get(ctx context.Context, q *dto.Query) ([]*dto.Node, error) {
	var nodes []*dto.Node
	mns, err := s.repo.Get(ctx, q)
	if err != nil {
		return nodes, err
	}
	for _, v := range mns {
		nodes = append(nodes, (&dto.Node{}).FillByModel(v))
	}
	return nodes, nil
}

func (s *NodeServiceImpl) Upsert(ctx context.Context, node *dto.Node) error {
	nodes, err := s.Get(ctx, &dto.Query{
		Query: "`id` = ?",
		Args:  []interface{}{node.ID},
	})
	if err != nil {
		return err
	}
	if len(nodes) > 0 {
		s.repo.Update(ctx, model.Node{ID: nodes[0].ID}, map[string]interface{}{
			"name":   node.Name,
			"addr":   node.Addr,
			"enable": node.Enable,
			"type":   node.Type,
		})
		return nil
	}
	n := model.Node{
		Name:     node.Name,
		Addr:     node.Addr,
		Type:     node.Type,
		Enable:   node.Enable,
		CreateAt: time.Now(),
	}
	return s.repo.Create(ctx, n)
}

func (s *NodeServiceImpl) Del(ctx context.Context, id int) error {
	return s.repo.Del(ctx, model.Node{ID: id})
}

func NewNodeService(repo repo.NodeRepo) *NodeServiceImpl {
	return &NodeServiceImpl{
		repo: repo,
	}
}
