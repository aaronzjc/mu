package dto

import (
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/pkg/helper"
)

type Favor struct {
	ID        int    `json:"id"`
	UserId    int    `json:"user_id"`
	Site      string `json:"site"`
	Key       string `json:"key"`
	OriginUrl string `json:"origin_url"`
	Title     string `json:"title"`
	CreateAt  string `json:"create_at"`
}

func (f *Favor) FillByModel(favor model.Favor) *Favor {
	f.ID = favor.ID
	f.UserId = favor.UserId
	f.Site = favor.Site
	f.OriginUrl = favor.OriginUrl
	f.Title = favor.Title
	f.CreateAt = helper.TimeToLocalStr(favor.CreateAt)
	return f
}

type FavorList struct {
	Tabs []*IndexSite `json:"tabs"`
	List []*Favor     `json:"list"`
}
