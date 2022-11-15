package site

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const SITE_TIEBA = "tieba"

var TiebaTabs = []SiteTab{
	{
		Tag:  "beiguo",
		Name: "抗压背锅吧",
		Url:  "https://tieba.baidu.com/f?ie=utf-8&kw=抗压背锅&fr=search",
	},
	{
		Tag:  "ruozhi",
		Name: "弱智吧",
		Url:  "https://tieba.baidu.com/f?ie=utf-8&kw=弱智&fr=search",
	},
}

type Tieba struct {
	Site
}

func (t *Tieba) GetSite() *Site {
	return &t.Site
}

func (t *Tieba) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range TiebaTabs {
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

func (t *Tieba) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := t.FetchData(link, nil, nil)
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
	reg := regexp.MustCompile(`.*/p/(\d+).*`)
	id := reg.ReplaceAllString(link, "$1")
	return id
}

func NewTieba() *Tieba {
	return &Tieba{
		Site{
			Name:     "贴吧",
			Key:      SITE_TIEBA,
			Root:     "https://tieba.baidu.com",
			Desc:     "鱼龙混杂的社区",
			CrawType: CrawHtml,
			Tabs:     TiebaTabs,
		},
	}
}

var _ Spider = &Tieba{}

func init() {
	RegistSite(SITE_TIEBA, NewTieba())
}
