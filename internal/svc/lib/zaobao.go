package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"mu/internal/util/tool"
	"time"
)

const SITE_ZAOBAO = "zaobao"

var ZaobaoTabs = []map[string]string{
	{
		"tag":  "focus",
		"url": "http://www.zaobao.com/",
		"name": "今日焦点",
	},
}

type Zaobao struct {
	Site
}

func (z *Zaobao) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range ZaobaoTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
			Sp:  z,
		}
		list = append(list, link)
	}

	return list, nil
}

func (z *Zaobao) CrawPage(link Link) (Page, error) {
	page, err := z.Craw(link, nil)
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
	return tool.MD55(link)
}
