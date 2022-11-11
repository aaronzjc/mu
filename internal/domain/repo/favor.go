package repo

import (
	"context"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
)

type FavorRepo interface {
	Get(context.Context, *dto.Query) ([]model.Favor, error)
	Create(context.Context, model.Favor) error
	Del(context.Context, model.Favor) error
	Sites(context.Context, *dto.Query) []string
}
