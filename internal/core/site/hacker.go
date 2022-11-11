package site

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aaronzjc/mu/pkg/helper"
)

const SITE_HACKER = "hacker"

var HackerTabs = []map[string]string{
	{
		"tag":  "new",
		"name": "最新",
		"url":  "https://news.ycombinator.com/",
	},
}

type Hacker struct {
	Site
}

var _ Spider = &Hacker{}

func init() {
	RegistSite(SITE_HACKER, &Hacker{
		Site{
			Name:     "Hacker",
			Key:      SITE_HACKER,
			Root:     "https://news.ycombinator.com/",
			Desc:     "Hacker News",
			CrawType: CrawHtml,
			Tabs:     HackerTabs,
		},
	})
}

func (h *Hacker) GetSite() *Site {
	return &h.Site
}

func (h *Hacker) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range HackerTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
		}
		list = append(list, link)
	}

	return list, nil
}

func (h *Hacker) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := h.FetchData(link, nil, headers)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find(".title").Find("a").Attr("href")
		text := s.Find(".title").Find(".titlelink").Text()
		if text == "" || url == "" {
			return
		}
		hot := Hot{
			Title:     text,
			OriginUrl: url,
		}
		hot.Key = h.FetchKey(hot.OriginUrl)
		if h.Key == "" {
			return
		}
		data = append(data, hot)
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (h *Hacker) FetchKey(link string) string {
	if link == "" {
		return ""
	}
	return helper.Md5(link)
}
