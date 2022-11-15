package site

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aaronzjc/mu/pkg/helper"
	"github.com/aaronzjc/mu/pkg/logger"
)

const SITE_CT = "chouti"

var ChoutiTabs = []SiteTab{
	{
		Url:  "/link/hot",
		Tag:  "hot",
		Name: "新热榜",
	},
	{
		Url:  "/top/24hr",
		Tag:  "24hr",
		Name: "24小时最热",
	},
	{
		Url:  "/top/72hr",
		Tag:  "72hr",
		Name: "3天最热",
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

func (c *Chouti) GetSite() *Site {
	return &c.Site
}

func (c *Chouti) BuildUrl() ([]Link, error) {
	var list []Link
	for _, item := range ChoutiTabs {
		link := Link{
			Key: item.Url,
			Url: fmt.Sprintf("%s%s", c.Root, item.Url),
			Tag: item.Tag,
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

func NewChouti() *Chouti {
	return &Chouti{
		Site{
			Name:     "抽屉",
			Key:      SITE_CT,
			Root:     "https://dig.chouti.com",
			Desc:     "抽屉新热榜",
			CrawType: CrawApi,
			Tabs:     ChoutiTabs,
		},
	}
}

var _ Spider = &Chouti{}

func init() {
	RegistSite(SITE_CT, NewChouti())
}
