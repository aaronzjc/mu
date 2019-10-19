package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"time"
)

const SITE_TIEBA = "tieba"

var TiebaTabs = []map[string]string{
	{
		"tag":  "lol",
		"name": "英雄联盟吧",
		"url":  "https://tieba.baidu.com/f?ie=utf-8&kw=英雄联盟&fr=search",
	},
	{
		"tag":  "kangya",
		"name": "抗压吧",
		"url":  "https://tieba.baidu.com/f?ie=utf-8&kw=抗压&fr=search",
	},
}

type Tieba struct {
	Site
}

func (t *Tieba) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range TiebaTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
			Sp:  t,
		}
		list = append(list, link)
	}

	return list, nil
}

func (t *Tieba) CrawPage(link Link) (Page, error) {
	page, err := t.Craw(link, nil)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".tl_shadow_new").Each(func(i int, s *goquery.Selection) {
		num := s.Find(".btn_icon").Text()
		url, _ := s.Find(".j_common").Attr("href")
		text := s.Find(".j_common").Find(".ti_title span").Text()
		if text == "" || url == "" {
			return
		}
		if num == "" {
			num = "0"
		}
		hot := Hot{
			Title:     fmt.Sprintf("%s - %s", text, num),
			OriginUrl: fmt.Sprintf("%s%s", t.Root, url),
		}
		hot.Key = t.FetchKey(hot.OriginUrl)
		if t.Key == "" {
			return
		}
		data = append(data, hot)
	})
	page.T = time.Now()
	page.List = data

	return page, nil
}

func (t *Tieba) FetchKey(link string) string {
	reg := regexp.MustCompile(".*/p/(\\d+).*")
	id := reg.ReplaceAllString(link, "$1")
	return id
}
