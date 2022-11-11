package store

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
	"github.com/aaronzjc/mu/internal/infra/cache"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type SiteRepoImpl struct {
	BaseRepoImpl
	cache *redis.Client
}

var _ repo.SiteRepo = &SiteRepoImpl{}

func (r *SiteRepoImpl) Get(ctx context.Context, q *dto.Query) ([]model.Site, error) {
	sites := []model.Site{}
	if err := r.prepare(q).Find(&sites).Error; err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return nil, errors.New("get sites err")
		}
	}
	return sites, nil
}

func (s *SiteRepoImpl) Create(ctx context.Context, site model.Site) error {
	return s.create(&site)
}

func (s *SiteRepoImpl) Update(ctx context.Context, site model.Site, data map[string]interface{}) error {
	return s.db.Model(&site).Updates(data).Error
}

func (s *SiteRepoImpl) Del(ctx context.Context, site model.Site) error {
	return s.db.Delete(&site).Error
}

func (s *SiteRepoImpl) GetNews(ctx context.Context, k string, kk string) (model.News, error) {
	data, err := s.cache.HGet(k, kk).Result()
	if err != nil {
		return model.News{}, errors.New("get news err")
	}
	var news model.News
	if err := json.Unmarshal([]byte(data), &news); err != nil {
		return model.News{}, errors.New("get news err")
	}
	return news, nil
}

func (s *SiteRepoImpl) SaveNews(ctx context.Context, site string, tag string, data string) error {
	redis := cache.Get()
	_, err := redis.HSet(site, tag, data).Result()
	return err
}

func NewSiteRepo() *SiteRepoImpl {
	base, _ := NewBaseImpl()
	cache := cache.Get()
	return &SiteRepoImpl{
		BaseRepoImpl: base,
		cache:        cache,
	}
}
