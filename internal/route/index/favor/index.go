package favor

import (
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/route/middleware"
	"mu/internal/svc/lib"
	"mu/internal/util/req"
	"time"
)

type ListForm struct {
	Site    string `form:"s"`
	Keyword string `form:"keyword"`
}

type AddForm struct {
	Key   string `json:"key"`
	Site  string `json:"site"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type RemoveForm struct {
	Site string `json:"site"`
	Key  string `json:"key"`
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
		UserId:    login.(int),
		Key:       r.Key,
		Site:      r.Site,
		OriginUrl: r.Url,
		Title:     r.Title,
		CreateAt:  time.Now(),
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
		req.JSON(c, req.CodeError, "添加失败", nil)
		return
	}

	req.JSON(c, req.CodeSuccess, "添加成功", nil)
	return
}

func Remove(c *gin.Context) {
	var err error
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
		Key:    r.Key,
		Site:   r.Site,
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
	var err error
	var r ListForm
	if err = c.ShouldBindQuery(&r); err != nil {
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

	q := model.Query{}

	var siteNames []string
	if keyword == "" {
		q = model.Query{
			Query: "`user_id` = ?",
			Args:  []interface{}{login.(int)},
		}
	} else {
		q = model.Query{
			Query: "`user_id` = ? AND `title` like ?",
			Args:  []interface{}{login.(int), "%" + keyword + "%"},
		}
	}
	siteNames = m.Config(q)
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
			"key":  site.Key,
			"tags": []string{},
		})
	}

	var favors []model.Favor
	if keyword == "" {
		q = model.Query{
			Query: "`site` = ? AND `user_id` = ?",
			Args:  []interface{}{site, login.(int)},
			Order: "`id` DESC",
		}
	} else {
		q = model.Query{
			Query: "`site` = ? AND `user_id` = ? AND `title` LIKE ?",
			Args:  []interface{}{site, login.(int), "%" + keyword + "%"},
			Order: "`id` DESC",
		}
	}
	favors, err = m.FetchRows(q)
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
