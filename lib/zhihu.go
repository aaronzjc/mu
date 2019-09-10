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
		"Cookie": `_xsrf=3cfddd0e-0c09-473d-a9f9-18e5fff1e1e8; _zap=2bb600bd-772d-4e8e-8ce4-acf8fb49ce92; cap_id="YmRlNDRkZTE0MTFlNDNiM2E2ZTExYTYxNzM2NzFiNDk=|1564627399|c4d052ccee2afa31e4d83f9d67d47dbbeb0bb9ac"; capsion_ticket="2|1:0|10:1567651271|14:capsion_ticket|44:ZWE2ZWMxODU3ZmM1NDQ0ODlhMzc0NzY2MmFkMTZkOGI=|bd04dc2f697f3307f2553eca6cd89fde948e8d6b6b000e4e20463603b2beb42e"; d_c0="ALDshH4Xgw-PTm5ahB1CV912PC6nOcTwOIM=|1559275457"; l_cap_id="MDVmNTI1NDBhMjFiNDRjODkwNzgyZmE3OTYyMTgyOTA=|1564627399|9c6fa2c2834cdbf07db470d8ffeaa39ddf1e4622"; q_c1=69697fefe2034642ae600963915231c7|1567651287000|1564627399000; r_cap_id="NmQ5NmIyMmU4NDg4NDIwMThjMDhjMmE0NTM1YmMzMjc=|1564627399|17755e09ab734667f12a5deccf99aadf2cc206bf"; tgw_l7_route=80f350dcd7c650b07bd7b485fcab5bf7; z_c0="2|1:0|10:1567651278|4:z_c0|92:Mi4xQ2pjSEFBQUFBQUFBc095RWZoZUREeVlBQUFCZ0FsVk56c05kWGdDdmNLRGZEZ0hvcmNpekFRUXVPTFBldmFiMXJR|a270dc45cfe084ac25d991bfb762a9b9c660cb657a109f7e72cf53f88c7f31e2"; tshl=; tst=r;`,
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
		T: page.T.Format("2006-01-02 15:04:05"),
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