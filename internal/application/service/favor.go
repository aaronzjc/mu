package service

import (
	"context"
	"errors"
	"time"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
)

type FavorService interface {
	UserFavors(context.Context, int, string, string) ([]*dto.Favor, error)
	UserFavorSites(context.Context, int, string) ([]string, error)

	Add(context.Context, *dto.Favor) error
	Del(context.Context, int, string, string) error
}

type FavorServiceImpl struct {
	repo repo.FavorRepo
}

func (s *FavorServiceImpl) Add(ctx context.Context, favor *dto.Favor) error {
	fs, _ := s.repo.Get(ctx, &dto.Query{
		Query: "`user_id` = ? AND `site` = ? AND `key` = ?",
		Args:  []interface{}{favor.UserId, favor.Site, favor.Key},
	})
	if len(fs) > 0 {
		return errors.New("重复内容")
	}
	if err := s.repo.Create(ctx, model.Favor{
		UserId:    favor.UserId,
		Key:       favor.Key,
		Site:      favor.Site,
		OriginUrl: favor.OriginUrl,
		Title:     favor.Title,
		CreateAt:  time.Now(),
	}); err != nil {
		return errors.New("添加失败")
	}
	return nil
}

func (s *FavorServiceImpl) Del(ctx context.Context, uid int, site string, key string) error {
	return s.repo.Del(ctx, model.Favor{
		UserId: uid,
		Site:   site,
		Key:    key,
	})
}

func (s *FavorServiceImpl) UserFavors(ctx context.Context, uid int, site string, keyword string) ([]*dto.Favor, error) {
	var favors []*dto.Favor
	q := &dto.Query{}
	if keyword == "" {
		q.Query = "`user_id` = ? AND `site` = ?"
		q.Args = []interface{}{uid, site}
	} else {
		q.Query = "`user_id` = ? AND `site` = ? AND `title` like ?"
		q.Args = []interface{}{uid, site, "%" + keyword + "%"}
	}

	mfs, _ := s.repo.Get(ctx, q)
	for _, v := range mfs {
		favors = append(favors, (&dto.Favor{}).FillByModel(v))
	}

	return favors, nil
}

func (s *FavorServiceImpl) UserFavorSites(ctx context.Context, uid int, keyword string) ([]string, error) {
	var sites []string
	q := &dto.Query{}
	if keyword == "" {
		q.Query = "`user_id` = ?"
		q.Args = []interface{}{uid}
	} else {
		q.Query = "`user_id` = ? AND `title` like ?"
		q.Args = []interface{}{uid, "%" + keyword + "%"}
	}

	sites = s.repo.Sites(ctx, q)
	return sites, nil
}

func NewFavorService(repo repo.FavorRepo) *FavorServiceImpl {
	return &FavorServiceImpl{
		repo: repo,
	}
}
