package api

import (
	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/gin-gonic/gin"
)

type Favor struct {
	svc service.FavorService
}

type ListForm struct {
	Site    string `form:"s"`
	Keyword string `form:"keyword"`
}

func (f *Favor) List(ctx *gin.Context) {
	var err error
	var r ListForm
	if err = ctx.ShouldBindQuery(&r); err != nil {
		Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}
	login, exist := ctx.Get(constant.LoginKey)
	if !exist {
		Resp(ctx, constant.CodeAuthFailed, "未登录", nil)
		return
	}

	favorList := dto.FavorList{
		Tabs: []string{},
		List: []*dto.Favor{},
	}

	sites, _ := f.svc.UserFavorSites(ctx, login.(int), r.Keyword)
	if len(sites) == 0 {
		Resp(ctx, constant.CodeSuccess, "success", favorList)
		return
	}
	favorList.Tabs = sites

	site := r.Site
	if site == "" {
		site = sites[0]
	}
	favorList.List, _ = f.svc.UserFavors(ctx, login.(int), site, r.Keyword)
	Resp(ctx, constant.CodeSuccess, "success", favorList)
}

type AddForm struct {
	Key   string `json:"key"`
	Site  string `json:"site"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

func (f *Favor) Add(ctx *gin.Context) {
	var err error
	var r AddForm
	if err = ctx.ShouldBindJSON(&r); err != nil {
		Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}

	login, exist := ctx.Get(constant.LoginKey)
	if !exist {
		Resp(ctx, constant.CodeForbidden, "未登录", nil)
		return
	}

	err = f.svc.Add(ctx, &dto.Favor{
		UserId:    login.(int),
		Key:       r.Key,
		Site:      r.Site,
		OriginUrl: r.Url,
		Title:     r.Title,
	})
	if err != nil {
		Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}
	Resp(ctx, constant.CodeSuccess, "success", nil)
}

type RemoveForm struct {
	Site string `json:"site"`
	Key  string `json:"key"`
}

func (f *Favor) Remove(ctx *gin.Context) {
	var err error
	var r RemoveForm
	if err = ctx.ShouldBindJSON(&r); err != nil {
		Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}

	login, exist := ctx.Get(constant.LoginKey)
	if !exist {
		Resp(ctx, constant.CodeForbidden, "未登录", nil)
		return
	}

	if err := f.svc.Del(ctx, login.(int), r.Site, r.Key); err != nil {
		Resp(ctx, constant.CodeError, "失败", nil)
		return
	}
	Resp(ctx, constant.CodeSuccess, "success", nil)
}

func NewFavor() *Favor {
	repo := store.NewFavorRepo()
	return &Favor{
		svc: service.NewFavorService(repo),
	}
}
