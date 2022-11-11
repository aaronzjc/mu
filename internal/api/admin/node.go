package admin

import (
	"strconv"

	"github.com/aaronzjc/mu/internal/api"
	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/gin-gonic/gin"
)

type Node struct {
	svc service.NodeService
}

type ListForm struct {
	Id      int    `form:"id"`
	Keyword string `form:"keyword"`
}

func (n *Node) List(ctx *gin.Context) {
	var r ListForm
	if err := ctx.ShouldBindQuery(&r); err != nil {
		api.Resp(ctx, constant.CodeError, "参数错误", nil)
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
	nodes, err := n.svc.Get(ctx, q)
	if err != nil {
		api.Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}

	api.Resp(ctx, constant.CodeSuccess, "success", nodes)
}

type UpsertForm struct {
	ID     int    `form:"id"`
	Name   string `form:"name"`
	Addr   string `form:"addr"`
	Type   int8   `form:"type"`
	Enable int8   `form:"enable"`
}

func (n *Node) Upsert(ctx *gin.Context) {
	var r UpsertForm
	if err := ctx.ShouldBind(&r); err != nil {
		api.Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}
	err = n.svc.Upsert(ctx, &dto.Node{
		ID:     id,
		Name:   r.Name,
		Addr:   r.Addr,
		Type:   r.Type,
		Enable: r.Enable,
	})
	if err != nil {
		api.Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}
	api.Resp(ctx, constant.CodeSuccess, "success", nil)
}

func (n *Node) Del(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.Resp(ctx, constant.CodeError, "参数错误", nil)
		return
	}
	err = n.svc.Del(ctx, id)
	if err != nil {
		api.Resp(ctx, constant.CodeError, err.Error(), nil)
		return
	}
	api.Resp(ctx, constant.CodeSuccess, "success", nil)
}

func NewNode() *Node {
	repo := store.NewNodeRepo()
	return &Node{
		svc: service.NewNodeService(repo),
	}
}
