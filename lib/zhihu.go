package lib

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

const RedisZhihu = "zhihu"

var ZhihuTabs = []map[string]string{
	{
		"tag":  "all",
		"url": "https://www.zhihu.com/billboard",
		"name": "知乎热榜",
	},
}

type Zhihu struct {
	Site
}

func (z *Zhihu) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range ZhihuTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
			Sp: z,
		}
		list = append(list, link)
	}

	return list, nil
}

func (z *Zhihu) CrawPage(link Link) (Page, error) {
	page, err := z.Craw(link)
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".HotList-item").Each(func(i int, s *goquery.Selection) {
		text := s.Find(".HotList-itemTitle").Text()
		if text == ""{
			return
		}

		data = append(data, Hot{
			Title: text,
			OriginUrl: fmt.Sprintf("https://www.zhihu.com/question/%s", ""),
		})
	})

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (v *Zhihu) Store(page Page) bool {
	hotJson := &HotJson{
		T: page.T,
		List: page.List,
	}

	data, err := json.Marshal(hotJson)
	if err != nil {
		log.Fatalf("[error] Json_encode zhihu error , err = %s\n", err.Error())
		return false
	}
	SaveToRedis(RedisZhihu, page.Link.Tag, string(data))

	log.Printf("[info] Store zhihu %s end", page.Link.Tag)

	return true
}