package hot

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/route/middleware"
	"mu/internal/svc/lib"
	"mu/internal/util/cache"
	"mu/internal/util/logger"
	"mu/internal/util/req"
)

type Tab struct {
	Name string      `json:"name"`
	Key  string      `json:"key"`
	Tags []model.Tag `json:"tags"`
}

type Item struct {
	lib.Hot
	Mark bool `json:"mark"`
}

func List(c *gin.Context) {
	client := cache.RedisConn()

	key := c.Request.URL.Query()["key"][0]
	hkey := c.Request.URL.Query()["hkey"][0]

	data, err := client.HGet(key, hkey).Result()

	if err != nil {
		logger.Info("aj req empty " + err.Error())
		req.JSON(c, req.CodeError, "没数据", nil)
		return
	}

	var hotJson lib.HotJson
	err = json.Unmarshal([]byte(data), &hotJson)
	if err != nil {
		req.JSON(c, req.CodeError, "请求列表失败", nil)
		return
	}

	// 遍历用户的收藏列表
	ckMap := make(map[string]bool)
	login, exist := c.Get(middleware.LoginUser)
	if exist {
		query := model.Query{
			Query: "`user_id` = ? AND `site` = ?",
			Args:  []interface{}{login.(int), key},
		}
		favors, err := (&model.Favor{}).FetchRows(query)
		if err != nil {
			req.JSON(c, req.CodeError, "请求收藏夹失败", nil)
			return
		}
		for _, v := range favors {
			ckMap[v.Key] = true
		}
	}

	result := make(map[string]interface{})
	var list []Item
	for _, val := range hotJson.List {
		exist := ckMap[val.Key]
		list = append(list, Item{
			val,
			exist,
		})
	}
	result["t"] = hotJson.T
	result["list"] = list

	logger.Info("req list key = " + key + " hkey = " + hkey)

	req.JSON(c, req.CodeSuccess, "成功", result)
	return
}

func Tabs(c *gin.Context) {
	var tabs []Tab
	var tags []model.Tag
	sites, _ := (&model.Site{}).FetchRows(model.Query{
		Query: "`enable` = ?",
		Args:  []interface{}{model.Enable},
	})
	for _, site := range sites {
		js, _ := site.FormatJson()
		tags = []model.Tag{}
		for _, tag := range js.Tags {
			if tag.Enable == 0 {
				continue
			}
			tags = append(tags, tag)
		}
		tabs = append(tabs, Tab{
			Name: site.Name,
			Key:  site.Key,
			Tags: tags,
		})
	}

	if len(tabs) == 0 {
		req.JSON(c, req.CodeSuccess, "成功", []struct{}{})
		return
	}

	req.JSON(c, req.CodeSuccess, "成功", tabs)
	return
}
