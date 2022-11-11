package site

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const SITE_V2EX = "v2ex"

var V2exTabs = []map[string]string{
	{
		"tag":  "all",
		"name": "全部",
	},
	{
		"tag":  "hot",
		"name": "最热",
	},
}

type V2ex struct {
	Site
}

var _ Spider = &V2ex{}

func init() {
	RegistSite(SITE_V2EX, &V2ex{
		Site{
			Name:     "v2ex",
			Key:      SITE_V2EX,
			Root:     "https://www.v2ex.com",
			Desc:     "way to explore",
			CrawType: CrawHtml,
			Tabs:     V2exTabs,
		},
	})
}

func (v *V2ex) GetSite() *Site {
	return &v.Site
}

func (v *V2ex) BuildUrl() ([]Link, error) {
	f := func(site string, tab string) string {
		return fmt.Sprintf("%s/?tab=%s", site, tab)
	}

	var list []Link
	for _, tab := range V2exTabs {
		url := f(v.Root, tab["tag"])
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
		}
		list = append(list, link)
	}

	return list, nil
}

func (v *V2ex) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := v.FetchData(link, nil, nil)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".cell tr").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find(".item_title").Find("a").Attr("href")
		text := s.Find(".item_title").Find("a").Text()
		comment := s.Find(".count_livid").Text()
		if text == "" || url == "" {
			return
		}
		if comment == "" {
			comment = "0"
		}
		h := Hot{
			Title:     text,
			OriginUrl: fmt.Sprintf("%s%s", v.Root, url),
			Rank: (func() float64 {
				val, _ := strconv.ParseFloat(comment, 32)
				return float64(val)
			})(),
		}
		h.Key = v.FetchKey(h.OriginUrl)
		if h.Key == "" {
			return
		}
		data = append(data, h)
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (v *V2ex) FetchKey(link string) string {
	reg := regexp.MustCompile(`.*/t/(\\d+).*`)
	id := reg.ReplaceAllString(link, "$1")
	return id
}
