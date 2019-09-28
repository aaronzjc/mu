package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"time"
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
			Sp:  v,
		}
		list = append(list, link)
	}

	return list, nil
}

func (v *V2ex) CrawPage(link Link) (Page, error) {
	page, err := v.Craw(link, nil)
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
		data = append(data, Hot{
			Title:     text,
			OriginUrl: fmt.Sprintf("%s%s", v.Root, url),
			Rank: (func() float64 {
				val, _ := strconv.ParseFloat(comment, 32)
				return float64(val)
			})(),
		})
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}