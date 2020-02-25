package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"mu/internal/util/tool"
	"strings"
	"time"
)

const SITE_GITHUB = "github"

var GithubTabs = []map[string]string{
	{
		"tag":  "trending",
		"url": "https://github.com/trending",
		"name": "Trending",
	},
	{
		"tag":  "trending-php",
		"url": "https://github.com/trending/php?since=daily",
		"name": "Trending-PHP",
	},
	{
		"tag":  "trending-go",
		"url": "https://github.com/trending/go?since=daily",
		"name": "Trending-Go",
	},
}

type Github struct {
	Site
}

func (g *Github) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range GithubTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
			Sp:  g,
		}
		list = append(list, link)
	}

	return list, nil
}

func (g *Github) CrawPage(link Link) (Page, error) {
	page, err := g.Craw(link, nil)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".Box .Box-row").Each(func(i int, s *goquery.Selection) {
		url := s.Find(" h1 a").AttrOr("href", "")
		desc := s.Find("p").Text()
		desc = strings.Trim(desc, "\n ")
		text := url[1:]

		if text == "" || url == "" {
			return
		}
		h := Hot{
			Title:     text + " - " + desc,
			OriginUrl: fmt.Sprintf("%s%s", g.Root, url),
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

func (g *Github) FetchKey(link string) string {
	if link == "" {
		return ""
	}
	return tool.MD55(link)
}
