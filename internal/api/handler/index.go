package handler

import (
	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/gin-gonic/gin"
)

type Index struct {
	svc service.SiteService
}

func (idx *Index) Sites(ctx *gin.Context) {
	sites, _ := idx.svc.ListOfIndex(ctx)
	Resp(ctx, constant.CodeSuccess, "", sites)
}

func (idx *Index) News(ctx *gin.Context) {
	k := ctx.Request.URL.Query()["key"][0]
	kk := ctx.Request.URL.Query()["hkey"][0]

	news, _ := idx.svc.News(ctx, k, kk)
	if news == nil {
		news = &dto.News{List: make([]dto.NewsItem, 0)}
	}
	Resp(ctx, constant.CodeSuccess, "success", news)
}

func NewIndex() *Index {
	repo := store.NewSiteRepo()
	favorRepo := store.NewFavorRepo()
	return &Index{
		svc: service.NewSiteService(repo, favorRepo),
	}
}
