package lib

import (
	"crawler/internal/util/cache"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
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

func (h *Hacker) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range HackerTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
			Sp:  h,
		}
		list = append(list, link)
	}

	return list, nil
}

func (h *Hacker) CrawPage(link Link) (Page, error) {
	page, err := h.Craw(link, nil)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find(".title").Find("a").Attr("href")
		text := s.Find(".title").Find(".storylink").Text()
		if text == "" || url == "" {
			return
		}
		data = append(data, Hot{
			Title:     text,
			OriginUrl: fmt.Sprintf("%s", url),
		})
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (h *Hacker) Store(page Page) bool {
	hotJson := &HotJson{
		T:    page.T.Format("2006-01-02 15:04:05"),
		List: page.List,
	}

	data, err := json.Marshal(hotJson)
	if err != nil {
		log.Printf("[error] Json_encode hacker news error , err = %s\n", err.Error())
		return false
	}
	cache.SaveToRedis(SITE_HACKER, page.Link.Tag, string(data))

	log.Printf("[info] Store hacker news %s end", page.Link.Tag)

	return true
}
