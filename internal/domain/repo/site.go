package repo

import (
	"context"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
)

type SiteRepo interface {
	Get(context.Context, *dto.Query) ([]model.Site, error)
	Create(context.Context, model.Site) error
	Update(context.Context, model.Site, map[string]interface{}) error
	Del(context.Context, model.Site) error
	GetNews(context.Context, string, string) (model.News, error)
	SaveNews(context.Context, string, string, string) error
}
