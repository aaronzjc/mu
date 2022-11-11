package site

import (
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const SITE_ZHIHU = "zhihu"

var ZhihuTabs = []map[string]string{
	{
		"tag":  "all",
		"url":  "https://www.zhihu.com/hot",
		"name": "知乎热榜",
	},
}

type Zhihu struct {
	Site
}

var _ Spider = &Zhihu{}

func init() {
	RegistSite(SITE_ZHIHU, &Zhihu{
		Site{
			Name:     "知乎",
			Key:      SITE_ZHIHU,
			Root:     "https://zhihu.com",
			Desc:     "知乎热榜",
			CrawType: CrawHtml,
			Tabs:     ZhihuTabs,
		},
	})
}

func (z *Zhihu) GetSite() *Site {
	return &z.Site
}

func (z *Zhihu) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range ZhihuTabs {
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

func (z *Zhihu) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := z.FetchData(link, nil, headers)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".HotList-list .HotItem-content").Each(func(i int, s *goquery.Selection) {
		url := s.Find("a").AttrOr("href", "")
		text := s.Find("h2").Text()
		if text == "" {
			return
		}
		hot := Hot{
			Title:     text,
			OriginUrl: url,
		}
		hot.Key = z.FetchKey(hot.OriginUrl)
		if hot.Key == "" {
			return
		}
		data = append(data, hot)
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (z *Zhihu) FetchKey(link string) string {
	reg := regexp.MustCompile(`.*/question/(\d+)`)
	id := reg.ReplaceAllString(link, "$1")
	return id
}
