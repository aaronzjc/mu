package lib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"time"
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

func (z *Zhihu) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range ZhihuTabs {
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

func (z *Zhihu) CrawPage(link Link) (Page, error) {
	page, err := z.Craw(link, map[string]string{
		"User-Agent": `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`,
		"Cookie": `_zap=f429a851-26f4-4483-8632-d3ae3b7ae654; d_c0="AIAow5kacQ6PTrqmH17zp-iAcfuUHny_NrM=|1540888443"; __gads=ID=129c4bb15876f84c:T=1544064126:S=ALNI_Ma2lau7-gf3DX_0TOj3Bn1KcWhahg; _xsrf=3s3cMDlLtV0TJAkq4pqCw3pkoQcPfSBF; __utmv=51854390.100--|2=registration_date=20130107=1^3=entry_date=20130107=1; _ga=GA1.2.1630895221.1567158654; __utma=51854390.1630895221.1567158654.1577428722.1578304905.3; __utmz=51854390.1578304905.3.3.utmcsr=mu.memosa.cn|utmccn=(referral)|utmcmd=referral|utmcct=/; q_c1=1156d4d4f7b5485da29c132f9136a3dd|1581998521000|1540888445000; _gid=GA1.2.1984428313.1583123523; Hm_lvt_98beee57fd2ef70ccdd5ca52b9740c49=1583123544,1583151604,1583151607,1583151665; tshl=; tst=r; Hm_lpvt_98beee57fd2ef70ccdd5ca52b9740c49=1583206612; _gat_gtag_UA_149949619_1=1; KLBRSID=d6f775bb0765885473b0cba3a5fa9c12|1583206613|1583205193; capsion_ticket="2|1:0|10:1583205814|14:capsion_ticket|44:NTA3MmViOGJhNWZiNGExODhkYmMxY2EyM2VkODc2NDY=|cbae5dba2997f5624c3b1cf5b373b49ba94fc2acd6271b64665aa75d2ef55c04"; z_c0="2|1:0|10:1583205821|4:z_c0|92:Mi4xQ2pjSEFBQUFBQUFBZ0NqRG1ScHhEaVlBQUFCZ0FsVk52UnRMWHdDOGZ3SG5fdTNDMHQxM3lIWUtXVFh2TjBURHdn|e4439cd63ac8533ceeef00ac24ac9d6eeaec79fde2d3b56a1eca5d96f2dac1ba"`,
	})
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
			OriginUrl: fmt.Sprintf("%s", url),
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
	reg := regexp.MustCompile(".*/question/(\\d+)")
	id := reg.ReplaceAllString(link, "$1")
	return id
}
