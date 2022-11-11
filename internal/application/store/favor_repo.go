package store

import (
	"context"
	"errors"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
)

type FavorRepoImpl struct {
	BaseRepoImpl
}

var _ repo.FavorRepo = &FavorRepoImpl{}

func (r *FavorRepoImpl) Get(ctx context.Context, q *dto.Query) ([]model.Favor, error) {
	favors := []model.Favor{}
	if err := r.prepare(q).Find(&favors).Error; err != nil {
		return nil, errors.New("get favors err")
	}
	return favors, nil
}

func (r *FavorRepoImpl) Create(ctx context.Context, favor model.Favor) error {
	return r.create(&favor)
}

func (r *FavorRepoImpl) Del(ctx context.Context, f model.Favor) error {
	r.db.Delete(&f)
	return nil
}

func (r *FavorRepoImpl) Sites(ctx context.Context, q *dto.Query) []string {
	var sites []string
	var favors []model.Favor
	if err := r.prepare(q).Select("DISTINCT(`site`)").Find(&favors).Error; err != nil {
		return sites
	}
	for _, v := range favors {
		sites = append(sites, v.Site)
	}
	return sites
}

func NewFavorRepo() *FavorRepoImpl {
	base, _ := NewBaseImpl()
	return &FavorRepoImpl{
		BaseRepoImpl: base,
	}
}
