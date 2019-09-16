package lib

import (
	"crawler/util/cache"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const SITE_CT = "chouti"

var ChoutiTabs = []map[string]string{
	{
		"url":  "/link/hot",
		"tag":  "hot",
		"name": "新热榜",
	},
	{
		"url":  "/top/24hr",
		"tag":  "24hr",
		"name": "24小时最热",
	},
	{
		"url":  "/top/72hr",
		"tag":  "72hr",
		"name": "3天最热最热",
	},
}

type HotList struct {
	Data    []map[string]interface{} `json:"data"`
	Code    int                      `json:"code"`
	Success bool                     `json:"success"`
}

type Chouti struct {
	Site
}

func (c *Chouti) BuildUrl() ([]Link, error) {
	var list []Link
	for _, item := range ChoutiTabs {
		link := Link{
			Key: item["url"],
			Url: fmt.Sprintf("%s%s", c.Root, item["url"]),
			Tag: item["tag"],
			Sp:  c,
		}

		list = append(list, link)
	}

	return list, nil
}

func (c *Chouti) CrawPage(link Link) (Page, error) {
	page, err := c.Craw(link, nil)

	var list HotList
	if err := json.Unmarshal([]byte(page.Content), &list); err != nil {
		log.Printf(err.Error())
		return Page{}, err
	}
	if err != nil {
		return Page{}, err
	}
	page.Json = list.Data

	var data []Hot
	for _, v := range page.Json {
		data = append(data, Hot{
			Id:        int(v["id"].(float64)),
			Title:     v["title"].(string),
			OriginUrl: v["originalUrl"].(string),
			Rank:      v["score"].(float64),
		})
	}

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (c *Chouti) Store(page Page) bool {
	hotJson := &HotJson{
		T:    page.T.Format("2006-01-02 15:04:05"),
		List: page.List,
	}

	data, err := json.Marshal(hotJson)
	if err != nil {
		log.Printf("[error] Json_encode chouti error , err = %s\n", err.Error())
		return false
	}
	cache.SaveToRedis(SITE_CT, page.Link.Tag, string(data))

	log.Printf("[info] Store chouti %s end", page.Link.Tag)

	return true
}
