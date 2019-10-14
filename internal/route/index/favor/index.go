package favor

import (
	"crawler/internal/model"
	"crawler/internal/route/middleware"
	"crawler/internal/svc/lib"
	"crawler/internal/util/req"
	"github.com/gin-gonic/gin"
	"time"
)

type ListForm struct {
	Site 	string 		`form:"s"`
	Keyword string 		`form:"keyword"`
}

type AddForm struct {
	Key		string 		`json:"key"`
	Site 	string 		`json:"site"`
	Url 	string 		`json:"url"`
	Title 	string 		`json:"title"`
}

type RemoveForm struct {
	Site 	string 		`json:"site"`
	Key  	string 		`json:"key"`
}

func Add(c *gin.Context) {
	var err error
	var r AddForm
	if err = c.ShouldBindJSON(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常 "+err.Error(), nil)
		return
	}

	login, exist := c.Get(middleware.LoginUser)
	if !exist {
		req.JSON(c, req.CodeForbidden, "未登录", nil)
		return
	}

	m := &model.Favor{
		UserId: login.(int),
		Key: r.Key,
		Site: r.Site,
		OriginUrl: r.Url,
		Title: r.Title,
		CreateAt: time.Now(),
	}

	exist, err = m.Exist()
	if err != nil {
		req.JSON(c, req.CodeError, "系统异常", nil)
		return
	}
	if exist {
		req.JSON(c, req.CodeError, "已经存在", nil)
		return
	}

	if err = m.Create(); err != nil {
		req.JSON(c, req.CodeError, "add failed " + err.Error(), nil)
		return
	}

	req.JSON(c, req.CodeSuccess, "add success ", nil)
	return
}

func Remove(c *gin.Context) {
	var err  error
	var r RemoveForm
	if err = c.ShouldBindJSON(&r); err != nil {
		req.JSON(c, req.CodeError, "参数错误", nil)
		return
	}

	login, exist := c.Get(middleware.LoginUser)
	if !exist {
		req.JSON(c, req.CodeForbidden, "未登录", nil)
		return
	}

	m := &model.Favor{
		UserId: login.(int),
		Key: r.Key,
		Site: r.Site,
	}

	exist, err = m.Exist()
	if err != nil {
		req.JSON(c, req.CodeError, "系统异常", nil)
		return
	}
	if !exist {
		req.JSON(c, req.CodeError, "不存在该记录", nil)
		return
	}

	if done := m.Del(); !done {
		req.JSON(c, req.CodeError, "delete failed", nil)
		return
	}

	req.JSON(c, req.CodeSuccess, "done", nil)
	return
}

func List(c *gin.Context) {
	var r ListForm
	if err := c.ShouldBindQuery(&r); err != nil {
		req.JSON(c, req.CodeError, "参数错误", nil)
		return
	}

	login, exist := c.Get(middleware.LoginUser)
	if !exist {
		req.JSON(c, req.CodeError, "未登录", nil)
		return
	}

	m := model.Favor{}

	site := r.Site
	keyword := r.Keyword

	var siteNames []string
	if keyword != "" {
		siteNames = m.Config(model.Query{
			Query: "`user_id` = ? AND `title` like ?",
			Args: []interface{}{login.(int), "%" + keyword + "%"},
		})
	} else {
		siteNames = m.Config(model.Query{})
	}
	if len(siteNames) == 0 {
		req.JSON(c, req.CodeSuccess, "成功", map[string]interface{}{
			"tabs": []string{},
			"list": []model.Favor{},
		})
		return
	}

	if site == "" {
		site = siteNames[0]
	}

	var tabs []map[string]interface{}
	for _, name := range siteNames {
		site := lib.NewSite(name)
		tabs = append(tabs, map[string]interface{}{
			"name": site.Name,
			"key": site.Key,
			"tags": []string{},
		})
	}

	var err error
	var favors []model.Favor
	if keyword == "" {
		favors, err = (&model.Favor{}).FetchRows(model.Query{
			Query: "`site` = ? AND `user_id` = ?",
			Args: []interface{}{site, login.(int)},
		})
	} else {
		favors, err = (&model.Favor{}).FetchRows(model.Query{
			Query: "`site` = ? AND `user_id` = ? AND `title` LIKE ?",
			Args: []interface{}{site, login.(int), "%" + keyword + "%"},
		})
	}
	if err != nil {
		req.JSON(c, req.CodeError, "获取失败", nil)
		return
	}

	var list []model.FavorJson
	for _, val := range favors {
		list = append(list, val.FormatJson())
	}

	req.JSON(c, req.CodeSuccess, "成功", map[string]interface{}{
		"tabs": tabs,
		"list": list,
	})
	return
}