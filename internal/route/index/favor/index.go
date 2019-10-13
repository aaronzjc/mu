package favor

import (
	"crawler/internal/model"
	"crawler/internal/route/middleware"
	"crawler/internal/svc/lib"
	"crawler/internal/util/req"
	"github.com/gin-gonic/gin"
	"time"
)

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

	if exist := m.Exist(); exist {
		req.JSON(c, req.CodeError, "already exist", nil)
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
	var r RemoveForm
	if err := c.ShouldBindJSON(&r); err != nil {
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

	if exist := m.Exist(); !exist {
		req.JSON(c, req.CodeError, "not exist", nil)
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
	site := c.Request.URL.Query()["s"][0]
	if site == "" {
		req.JSON(c, req.CodeError, "参数错误", nil)
		return
	}

	login, exist := c.Get(middleware.LoginUser)
	if !exist {
		req.JSON(c, req.CodeError, "未登录", nil)
		return
	}

	favors, err := (&model.Favor{}).FetchRows("`site` = ? AND `user_id` = ?", site, login.(int))
	if err != nil {
		req.JSON(c, req.CodeError, "获取失败", nil)
		return
	}

	req.JSON(c, req.CodeSuccess, "成功", favors)
	return
}

func Config(c *gin.Context) {
	sitesNames := (&model.Favor{}).Config()

	var result []map[string]interface{}
	for _, name := range sitesNames {
		site := lib.NewSite(name)
		result = append(result, map[string]interface{}{
			"name": site.Name,
			"key": site.Key,
			"tags": []string{},
		})
	}
	req.JSON(c, req.CodeSuccess, "成功", result)
	return
}