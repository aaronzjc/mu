package front

import (
	"crawler/internal/svc/lib"
	"crawler/internal/util/cache"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Tag struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Tab struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Tags []Tag  `json:"tags"`
}

func Aj(c *gin.Context) {
	client := cache.RedisConn()
	defer client.Close()

	key := c.Request.URL.Query()["key"][0]
	hkey := c.Request.URL.Query()["hkey"][0]

	data, err := client.HGet(key, hkey).Result()

	if err != nil {
		log.Println("[info] aj req empty " + err.Error())
		c.JSON(http.StatusOK, []byte(`{"list": [], "t":""}`))
		return
	}

	var hotJson lib.HotJson
	err = json.Unmarshal([]byte(data), &hotJson)
	if err != nil {
		log.Println("[error] aj req error " + err.Error())
		c.JSON(http.StatusOK, []byte(`{"list": [], "t":""}`))
		return
	}

	js, _ := json.Marshal(hotJson)

	c.JSON(http.StatusOK, []byte(js))
}

func Config(c *gin.Context) {
	var tabs []Tab

	var fetchTags = func(tabs []map[string]string) []Tag {
		var tags []Tag
		for _, v := range tabs {
			tags = append(tags, Tag{
				Name: v["name"],
				Key:  v["tag"],
			})
		}

		return tags
	}

	sites := []string{
		lib.SITE_V2EX,
		lib.SITE_CT,
		lib.SITE_ZHIHU,
		lib.SITE_WEIBO,
		lib.SITE_HACKER,
	}

	for _, s := range sites {
		st := lib.NewSite(s)
		tabs = append(tabs, Tab{
			Name: st.Name,
			Key:  st.Key,
			Tags: fetchTags(st.Tabs),
		})
	}

	data, _ := json.Marshal(tabs)

	c.JSON(http.StatusOK, data)
}
