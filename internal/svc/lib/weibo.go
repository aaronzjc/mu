package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"mu/internal/util/tool"
	"time"
)

const SITE_WEIBO = "weibo"

var WeiboTabs = []map[string]string{
	{
		"tag":  "hot",
		"url":  "https://s.weibo.com/top/summary?cate=realtimehot",
		"name": "热搜",
	},
}

type Weibo struct {
	Site
}

func (w *Weibo) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range WeiboTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
			Sp:  w,
		}
		list = append(list, link)
	}

	return list, nil
}

func (w *Weibo) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := w.FetchData(link, nil, headers)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find("tbody td.td-02").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a").First()
		text, url := link.Text(), link.AttrOr("href", "#")
		if text == "" {
			return
		}
		hot := Hot{
			Title:     text,
			OriginUrl: fmt.Sprintf("%s%s", w.Root, url),
			Rank: 0,
		}
		hot.Key = w.FetchKey(hot.OriginUrl)
		if hot.Key == "" {
			return
		}
		data = append(data, hot)
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (w *Weibo) FetchKey(link string) string {
	if link == "" {
		return ""
	}
	return tool.MD55(link)
}
