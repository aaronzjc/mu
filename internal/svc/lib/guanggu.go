package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"time"
)

const SITE_GUANGGU = "guanggu"

var GuangGuTabs = []map[string]string{
	{
		"tag":  "default",
		"name": "默认",
		"url": "http://www.guanggoo.com/",
	},
	{
		"tag":  "latest",
		"name": "最新",
		"url": "http://www.guanggoo.com/?tab=latest",
	},
}

type Guanggu struct {
	Site
}

func (g *Guanggu) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range GuangGuTabs {
		link := Link{
			Key: tab["url"],
			Url: tab["url"],
			Tag: tab["tag"],
			Sp:  g,
		}
		list = append(list, link)
	}

	return list, nil
}

func (g *Guanggu) CrawPage(link Link) (Page, error) {
	page, err := g.Craw(link, nil)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".topic-item").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find(".main .title").Find("a").Attr("href")
		text := s.Find(".main .title").Find("a").Text()
		comment := s.Find(".count").Find("a").Text()
		if text == "" || url == "" {
			return
		}
		if comment == "" {
			comment = "0"
		}
		h := Hot{
			Title:     text,
			OriginUrl: fmt.Sprintf("%s%s", g.Root, url),
			Rank: (func() float64 {
				val, _ := strconv.ParseFloat(comment, 32)
				return float64(val)
			})(),
		}
		h.Key = g.FetchKey(h.OriginUrl)
		if h.Key == "" {
			return
		}
		data = append(data, h)
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (g *Guanggu) FetchKey(link string) string {
	reg := regexp.MustCompile(".*/t/(\\d+).*")
	id := reg.ReplaceAllString(link, "$1")
	return id
}
