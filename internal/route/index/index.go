package index

import (
	"crawler/internal/model"
	"crawler/internal/svc/lib"
	"crawler/internal/util/cache"
	"crawler/internal/util/req"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

type Tab struct {
	Name string      `json:"name"`
	Key  string      `json:"key"`
	Tags []model.Tag `json:"tags"`
}

func Aj(c *gin.Context) {
	client := cache.RedisConn()
	defer client.Close()

	key := c.Request.URL.Query()["key"][0]
	hkey := c.Request.URL.Query()["hkey"][0]

	data, err := client.HGet(key, hkey).Result()

	if err != nil {
		log.Println("[info] aj req empty " + err.Error())
		req.JSON(c, req.CodeError, "没数据", nil)
		return
	}

	var hotJson lib.HotJson
	err = json.Unmarshal([]byte(data), &hotJson)
	if err != nil {
		log.Println("[error] aj req error " + err.Error())
		req.JSON(c, req.CodeError, "请求失败", nil)
		return
	}

	req.JSON(c, req.CodeSuccess, "没数据", hotJson)
	return
}

func Config(c *gin.Context) {
	var tabs []Tab
	sites, _ := (&model.Site{}).FetchRows("`enable` = ?", model.Enable)
	for _, site := range sites {
		js, _ := site.FormatJson()
		tabs = append(tabs, Tab{
			Name: site.Name,
			Key:  site.Key,
			Tags: js.Tags,
		})
	}

	if len(tabs) == 0 {
		req.JSON(c, req.CodeSuccess, "成功", []struct{}{})
		return
	}

	req.JSON(c, req.CodeSuccess, "成功", tabs)
	return
}
