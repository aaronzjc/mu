package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
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

func (w *Weibo) CrawPage(link Link) (Page, error) {
	page, err := w.Craw(link, nil)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".list_a li").Each(func(i int, s *goquery.Selection) {
		text := s.Find("a").Find("span").Text()
		re := regexp.MustCompile(`\d+\s$`)
		text = re.ReplaceAllString(text, "")
		url := s.Find("a").AttrOr("href", "#")
		rank := s.Find("a").Find("span").Find("em").Text()
		if text == "" {
			return
		}

		data = append(data, Hot{
			Title:     text,
			OriginUrl: fmt.Sprintf("%s%s", w.Root, url),
			Rank: (func() float64 {
				f, _ := strconv.ParseFloat(rank, 64)
				return f
			})(),
		})
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}
