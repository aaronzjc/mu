package node

import (
	"crawler/internal/model"
	"crawler/internal/util/req"
	"github.com/gin-gonic/gin"
	"time"
)

type InfoForm struct {
	Id 			int 		`form:"id" binding:"required"`
}

type ListForm struct {
	Keyword 	string 		`form:"keyword"`
}

type UpsertForm struct {
	ID 			int 		`form:"id"`
	Name 		string		`form:"name"`
	Addr 			string 		`form:"addr"`
	Type 		int8 		`form:"type"`
	Enable 		int8 		`form:"enable"`
}

func Info(c *gin.Context) {
	var r InfoForm
	if err := c.ShouldBindQuery(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	m := &model.Node{
		ID: r.Id,
	}

	node, err := m.FetchInfo()
	if err != nil || node.ID <= 0 {
		req.JSON(c, req.CodeError, "失败", nil)
		return
	}

	json, _ := node.FormatJson()

	req.JSON(c, req.CodeSuccess, "成功", json)
	return
}

func List(c *gin.Context) {
	var r ListForm
	var err error
	if err = c.ShouldBind(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	m := &model.Node{}
	var nodes []model.Node
	if r.Keyword != "" {
		nodes, err = m.FetchRows("name like ?", "%" + r.Keyword + "%")
	} else {
		nodes, err = m.FetchRows("1=1")
	}
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}

	var result []model.NodeJson
	for _, node := range nodes {
		item, _ := node.FormatJson()
		result = append(result, item)
	}

	req.JSON(c, req.CodeSuccess, "成功", result)
	return
}

func CreateOrUpdateNode(c *gin.Context) {
	var err error
	var r UpsertForm
	if err = c.ShouldBind(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常 " + err.Error(), nil)
		return
	}

	m := &model.Node{
		Name: r.Name,
		Addr: r.Addr,
		Type: r.Type,
		Enable: r.Enable,
	}
	err = m.CheckArgs()
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}

	if r.ID > 0 {
		data := make(map[string]interface{})
		data["name"] = r.Name
		data["addr"] = r.Addr
		data["enable"] = r.Enable
		data["type"] = r.Type

		m.ID = r.ID
		err = m.Update(data)
		if err != nil {
			req.JSON(c, req.CodeError, "更新节点失败", nil)
			return
		}
		req.JSON(c, req.CodeSuccess, "成功", nil)
		return
	} else {
		m.CreateAt = time.Now()
		err := m.Create()
		if err != nil {
			req.JSON(c, req.CodeError, "插入节点失败", nil)
			return
		}
		req.JSON(c, req.CodeSuccess, "成功", nil)
		return
	}
}

func Del(c *gin.Context) {
	var r InfoForm
	if err := c.ShouldBindQuery(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	m := &model.Node{
		ID: r.Id,
	}

	if ok := m.Del(); !ok {
		req.JSON(c, req.CodeError, "删除失败", nil)
		return
	}

	// 删除了节点，需要更新站点配置的节点
	(&model.Site{}).FixNodeId(r.Id)

	req.JSON(c, req.CodeSuccess, "删除成功", nil)
	return
}