package site

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"mu/internal/model"
	"mu/internal/svc/rpc"
	"mu/internal/util/config"
	"mu/internal/util/logger"
	"mu/internal/util/req"
	"time"
)

type InfoForm struct {
	Id int `form:"id" binding:"required"`
}

type ListForm struct {
	Keyword string `form:"keyword"`
}

type UpdateForm model.SiteJson

func Info(c *gin.Context) {
	var r InfoForm
	if err := c.ShouldBindQuery(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	m := &model.Site{
		ID: r.Id,
	}

	site, err := m.FetchInfo()
	if err != nil || site.ID <= 0 {
		req.JSON(c, req.CodeError, "失败", nil)
		return
	}

	jsonObj, err := site.FormatJson()
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}
	req.JSON(c, req.CodeSuccess, "成功", jsonObj)
	return
}

func List(c *gin.Context) {
	var r ListForm
	var err error
	if err = c.ShouldBind(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	m := &model.Site{}
	var sites []model.Site
	if r.Keyword != "" {
		sites, err = m.FetchRows(model.Query{
			Query: "name like ?",
			Args:  []interface{}{"%" + r.Keyword + "%"},
		})
	} else {
		sites, err = m.FetchRows(model.Query{})
	}
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}

	var result []model.SiteJson
	for _, site := range sites {
		item, err := site.FormatJson()
		if err != nil {
			req.JSON(c, req.CodeError, err.Error(), nil)
			return
		}
		result = append(result, item)
	}

	nodes, _ := (&model.Node{}).FetchRows(model.Query{
		Query: "`enable` = ?",
		Args:  []interface{}{model.Enable},
	})
	nodeJson := make(map[int]model.Node)
	for _, node := range nodes {
		n, _ := node.FormatJson()
		nodeJson[node.ID] = n
	}

	req.JSON(c, req.CodeSuccess, "成功", map[string]interface{}{
		"nodeList": nodeJson,
		"siteList": result,
	})
	return
}

func UpdateSite(c *gin.Context) {
	var err error
	var r UpdateForm
	if err = c.ShouldBindJSON(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常 "+err.Error(), nil)
		return
	}

	tagBytes, _ := json.Marshal(r.Tags)
	hostsBytes, _ := json.Marshal(r.NodeHosts)
	reqHeaders, _ := json.Marshal(r.ReqHeaders)

	m := model.Site{
		ID:         r.ID,
		Name:       r.Name,
		Key:        r.Key,
		Desc:       r.Desc,
		Cron:       r.Cron,
		Tags:       string(tagBytes),
		Enable:     r.Enable,
		NodeOption: r.NodeOption,
		NodeType:   r.NodeType,
		NodeHosts:  string(hostsBytes),
		ReqHeaders: string(reqHeaders),
	}
	err = m.CheckArgs()
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}

	data := make(map[string]interface{})
	data["name"] = m.Name
	data["key"] = m.Key
	data["desc"] = m.Desc
	data["cron"] = m.Cron
	data["tags"] = m.Tags
	data["enable"] = m.Enable
	data["node_option"] = m.NodeOption
	data["node_type"] = m.NodeType
	data["node_hosts"] = m.NodeHosts
	data["req_headers"] = m.ReqHeaders

	om, _ := (&model.Site{}).FetchRow(model.Query{
		Query: "`id` = ?",
		Args:  []interface{}{m.ID},
	})

	err = m.Update(data)
	if err != nil {
		req.JSON(c, req.CodeError, "更新站点失败", nil)
		return
	}

	// 更新调度任务
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(config.NewConfig().Commander.Addr, opts...)
	if err != nil {
		logger.Error("connect error " + err.Error())
	} else {
		client := rpc.NewCommanderClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		client.UpdateCron(ctx, &rpc.Cron{Site: om.Key})
		logger.Info("remote update cron [%s]", om.Key)
	}

	req.JSON(c, req.CodeSuccess, "成功", nil)
	return
}

func Debug(c *gin.Context) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(config.NewConfig().Commander.Addr, opts...)
	if err != nil {
		logger.Error("connect error " + err.Error())
	} else {
		client := rpc.NewCommanderClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		res, err := client.Debug(ctx, &rpc.Empty{})
		if err != nil {
			logger.Error("Debug error " + err.Error())
			c.String(req.CodeError, "错误")
		} else {
			c.String(req.CodeSuccess, res.Res)
		}
	}
	return
}
