package admin

import (
	"strconv"

	"github.com/aaronzjc/mu/internal/api/handler"
	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type Site struct {
	svc     service.SiteService
	nodeSvc service.NodeService
	crawSvc service.CrawService
}

type SiteQueryForm struct {
	Id      int    `form:"id"`
	Keyword string `form:"keyword"`
}

func (s *Site) List(ctx *gin.Context) {
	var r SiteQueryForm
	if err := ctx.ShouldBindQuery(&r); err != nil {
		handler.Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}

	q := &dto.Query{}
	if r.Id > 0 {
		q.Query = "`id` = ?"
		q.Args = []interface{}{r.Id}
	}
	if r.Keyword != "" {
		q.Query = "`title` = ?"
		q.Args = []interface{}{r.Keyword}
	}
	sites, err := s.svc.Get(ctx, q)
	if err != nil {
		handler.Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}
	nodes, _ := s.nodeSvc.Get(ctx, &dto.Query{
		Query: "`enable` = ?",
		Args:  []interface{}{model.Enable},
	})
	nodeList := make(map[int]*dto.Node)
	for _, v := range nodes {
		nodeList[v.ID] = v
	}
	handler.Resp(ctx, constant.CodeSuccess, "success", map[string]interface{}{
		"nodeList": nodeList,
		"siteList": sites,
	})
}

type UpsertSiteForm struct {
	ID     int    `form:"id"`
	Name   string `form:"name"`
	Addr   string `form:"addr"`
	Type   int8   `form:"type"`
	Enable int8   `form:"enable"`
}

func (s *Site) Upsert(ctx *gin.Context) {
	var r dto.Site
	if err := ctx.ShouldBind(&r); err != nil {
		handler.Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handler.Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}
	r.ID = id
	err = s.svc.Upsert(ctx, &r)
	if err != nil {
		handler.Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}

	handler.Resp(ctx, constant.CodeSuccess, "success", nil)
}

func (s *Site) Del(ctx *gin.Context) {
	var r SiteQueryForm
	if err := ctx.ShouldBindQuery(&r); err != nil {
		handler.Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}
	err := s.svc.Del(ctx, r.Id)
	if err != nil {
		handler.Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}
	handler.Resp(ctx, constant.CodeSuccess, "success", nil)
}

func (s *Site) Craw(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if id <= 0 {
		handler.Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}
	sites, _ := s.svc.Get(ctx, &dto.Query{
		Query: "`id` = ?",
		Args:  []interface{}{id},
	})
	if len(sites) <= 0 {
		handler.Resp(ctx, constant.CodeError, "站点不存在", nil)
		return
	}
	err := s.crawSvc.Craw(ctx, sites[0])
	if err != nil {
		handler.Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}
	handler.Resp(ctx, constant.CodeSuccess, "success", nil)
}

func NewSite() *Site {
	return &Site{
		svc:     service.NewSiteService(store.NewSiteRepo(), nil),
		nodeSvc: service.NewNodeService(store.NewNodeRepo()),
		crawSvc: service.NewCrawService(),
	}
}
