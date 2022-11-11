package store

import (
	"context"
	"errors"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
)

type NodeRepoImpl struct {
	BaseRepoImpl
}

var _ repo.NodeRepo = &NodeRepoImpl{}

func (r *NodeRepoImpl) Get(ctx context.Context, q *dto.Query) ([]model.Node, error) {
	nodes := []model.Node{}
	if err := r.prepare(q).Find(&nodes).Error; err != nil {
		return nil, errors.New("get nodes err")
	}
	return nodes, nil
}

func (r *NodeRepoImpl) Create(ctx context.Context, node model.Node) error {
	return r.create(&node)
}

func (r *NodeRepoImpl) Update(ctx context.Context, node model.Node, data map[string]interface{}) error {
	return r.db.Model(&node).Updates(data).Error
}

func (r *NodeRepoImpl) Del(ctx context.Context, node model.Node) error {
	return r.db.Delete(&node).Error
}

func NewNodeRepo() *NodeRepoImpl {
	base, _ := NewBaseImpl()
	return &NodeRepoImpl{BaseRepoImpl: base}
}
