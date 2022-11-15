package site

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aaronzjc/mu/pkg/helper"
)

const SITE_GITHUB = "github"

var GithubTabs = []SiteTab{
	{
		Tag:  "trending",
		Url:  "https://github.com/trending",
		Name: "Trending",
	},
	{
		Tag:  "trending-php",
		Url:  "https://github.com/trending/php?since=daily",
		Name: "Trending-PHP",
	},
	{
		Tag:  "trending-go",
		Url:  "https://github.com/trending/go?since=daily",
		Name: "Trending-Go",
	},
}

type Github struct {
	Site
}

func (g *Github) GetSite() *Site {
	return &g.Site
}

func (g *Github) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range GithubTabs {
		url := tab.Url
		link := Link{
			Key: url,
			Url: url,
			Tag: tab.Tag,
		}
		list = append(list, link)
	}

	return list, nil
}

func (g *Github) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := g.FetchData(link, nil, nil)
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
			Title:     strings.Replace(text, "/", " • ", 1),
			Desc:      desc,
			OriginUrl: fmt.Sprintf("%s%s", g.Root, url),
			Card:      CardRichText,
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
	return helper.Md5(link)
}

func NewGithub() *Github {
	return &Github{
		Site{
			Name:     "Github",
			Key:      SITE_GITHUB,
			Root:     "https://github.com",
			Desc:     "Github.com",
			CrawType: CrawHtml,
			Tabs:     GithubTabs,
		},
	}
}

var _ Spider = &Github{}

func init() {
	RegistSite(SITE_GITHUB, NewGithub())
}
