package repo

import (
	"context"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
)

type NodeRepo interface {
	Get(context.Context, *dto.Query) ([]model.Node, error)
	Create(context.Context, model.Node) error
	Update(context.Context, model.Node, map[string]interface{}) error
	Del(context.Context, model.Node) error
}
