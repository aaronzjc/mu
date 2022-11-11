package site

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aaronzjc/mu/pkg/helper"
	"github.com/aaronzjc/mu/pkg/logger"
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
		"name": "3天最热",
	},
}

type ChoutiList struct {
	Data    []map[string]interface{} `json:"data"`
	Code    int                      `json:"code"`
	Success bool                     `json:"success"`
}

type Chouti struct {
	Site
}

var _ Spider = &Chouti{}

func init() {
	RegistSite(SITE_CT, &Chouti{
		Site{
			Name:     "抽屉",
			Key:      SITE_CT,
			Root:     "https://dig.chouti.com",
			Desc:     "抽屉新热榜",
			CrawType: CrawApi,
			Tabs:     ChoutiTabs,
		},
	})
}

func (c *Chouti) GetSite() *Site {
	return &c.Site
}

func (c *Chouti) BuildUrl() ([]Link, error) {
	var list []Link
	for _, item := range ChoutiTabs {
		link := Link{
			Key: item["url"],
			Url: fmt.Sprintf("%s%s", c.Root, item["url"]),
			Tag: item["tag"],
		}

		list = append(list, link)
	}

	return list, nil
}

func (c *Chouti) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := c.FetchData(link, nil, nil)
	if err != nil {
		return Page{}, err
	}

	var list ChoutiList
	if err := json.Unmarshal([]byte(page.Content), &list); err != nil {
		logger.Error(err.Error())
		return Page{}, err
	}
	page.Json = list.Data

	var data []Hot
	for _, v := range page.Json {
		h := Hot{
			Title:     v["title"].(string),
			OriginUrl: v["originalUrl"].(string),
			Rank:      v["score"].(float64),
		}
		h.Key = c.FetchKey(h.OriginUrl)
		if h.Key == "" {
			continue
		}
		data = append(data, h)
	}

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (c *Chouti) FetchKey(link string) string {
	if link == "" {
		return ""
	}
	return helper.Md5(link)
}
