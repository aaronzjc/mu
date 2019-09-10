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
		"url": "https://www.zhihu.com/hot",
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

	page, err := z.Craw(link, map[string]string{
		"User-Agent": `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`,
		"Cookie": `_zap=09ee8132-fd2b-43d3-9562-9d53a41a4ef5; d_c0="AGDv-acVoQ-PTvS01pG8OiR9v_9niR11ukg=|1561288241"; capsion_ticket="2|1:0|10:1561288248|14:capsion_ticket|44:NjE1ZTMxMjcxYjlhNGJkMjk5OGU4NTRlNDdkZTJhNzk=|7aefc35b3dfd27b74a087dd1d15e7a6bb9bf5c6cdbe8471bc20008feb67e7a9f"; z_c0="2|1:0|10:1561288250|4:z_c0|92:Mi4xeGZsekFBQUFBQUFBWU9fNXB4V2hEeVlBQUFCZ0FsVk5PcXo4WFFBNWFFRnhYX2h0ZFZpWTQ5T3dDMGh5ZTV1bjB3|0cee5ae41ff7053a1e39d96df2450077d37cc9924b337584cf006028b0a02f30"; q_c1=ae65e92b2bbf49e58dee5b2b29e1ffb3|1561288383000|1561288383000; tgw_l7_route=f2979fdd289e2265b2f12e4f4a478330; _xsrf=f8139fd6-b026-4f01-b860-fe219aa63543; tst=h; tshl=`,
	})
	if err != nil {
		return Page{}, err
	}
	var data []Hot
	doc := page.Doc
	doc.Find(".HotList-list .HotItem-content").Each(func(i int, s *goquery.Selection) {
		url := s.Find("a").AttrOr("href", "")
		text := s.Find("h2").Text()
		if text == ""{
			return
		}

		data = append(data, Hot{
			Title: text,
			OriginUrl: fmt.Sprintf("%s", url),
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
		log.Printf("[error] Json_encode zhihu error , err = %s\n", err.Error())
		return false
	}
	SaveToRedis(RedisZhihu, page.Link.Tag, string(data))

	log.Printf("[info] Store zhihu %s end", page.Link.Tag)

	return true
}