package site

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aaronzjc/mu/pkg/helper"
)

const SITE_ZAOBAO = "zaobao"

var ZaobaoTabs = []map[string]string{
	{
		"tag":  "focus",
		"url":  "http://www.zaobao.com/",
		"name": "今日焦点",
	},
}

type Zaobao struct {
	Site
}

var _ Spider = &Zaobao{}

func init() {
	RegistSite(SITE_ZAOBAO, &Zaobao{
		Site{
			Name:     "联合早报",
			Key:      SITE_ZAOBAO,
			Root:     "http://www.zaobao.com",
			Desc:     "新加坡新闻",
			CrawType: CrawHtml,
			Tabs:     ZaobaoTabs,
		},
	})
}

func (z *Zaobao) GetSite() *Site {
	return &z.Site
}

func (z *Zaobao) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range ZaobaoTabs {
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

func (z *Zaobao) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := z.FetchData(link, nil, nil)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find("#piping-hot .post-item-special p a").Each(func(i int, s *goquery.Selection) {
		url := s.AttrOr("href", "")
		text := s.Text()
		if text == "" || url == "" {
			return
		}
		h := Hot{
			Title:     text,
			OriginUrl: fmt.Sprintf("%s%s", z.Root, url),
		}
		h.Key = z.FetchKey(h.OriginUrl)
		if h.Key == "" {
			return
		}
		data = append(data, h)
	})
	doc.Find("#piping-hot a").Each(func(i int, s *goquery.Selection) {
		url := s.AttrOr("href", "")
		text := s.Find("span.post-title").Text()
		if text == "" || url == "" {
			return
		}
		h := Hot{
			Title:     text,
			OriginUrl: fmt.Sprintf("%s%s", z.Root, url),
		}
		h.Key = z.FetchKey(h.OriginUrl)
		if h.Key == "" {
			return
		}
		data = append(data, h)
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (z *Zaobao) FetchKey(link string) string {
	if link == "" {
		return ""
	}
	return helper.Md5(link)
}
