package site

import (
	"encoding/json"
	"fmt"
	"time"
)

const SITE_WBVIDEO = "wbvideo"

var WbvideoTabs = []SiteTab{
	{
		Tag:  "all",
		Url:  "https://weibo.com/tv/api/component?page=%2Ftv%2Fbillboard",
		Name: "全站",
		Args: map[string]string{
			"cid": "4418213501411061",
		},
	},
	{
		Tag:  "funny",
		Url:  "https://weibo.com/tv/api/component?page=%2Ftv%2Fbillboard%2F4418219809678869",
		Name: "搞笑幽默",
		Args: map[string]string{
			"cid": "4418219809678869",
		},
	},
}

type WbVideoList struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Videos struct {
			Next int `json:"next_cursor"`
			List []struct {
				Title     string `json:"title"`
				Cover     string `json:"cover_image"`
				Id        int64  `json:"mid"`
				Oid       string `json:"oid"`
				Date      string `json:"date"`
				PlayCount string `json:"play_count"`
			} `json:"list"`
		} `json:"Component_Billboard_Billboardlist"`
	} `json:"data"`
}

type Wbvideo struct {
	Site
}

func (w *Wbvideo) GetSite() *Site {
	return &w.Site
}

func (w *Wbvideo) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range WbvideoTabs {
		url := tab.Url
		link := Link{
			Key:    tab.Args["cid"],
			Url:    url,
			Tag:    tab.Tag,
			Method: "POST",
		}
		list = append(list, link)
	}

	return list, nil
}

func (w *Wbvideo) CrawPage(link Link, headers map[string]string) (res Page, err error) {
	var page Page
	var hotList []Hot
	var nextCursor int
	post := make(map[string]map[string]interface{})
	for {
		if nextCursor == 0 {
			post = map[string]map[string]interface{}{
				"Component_Billboard_Billboardcategory": {},
				"Component_Billboard_Billboardlist": {
					"cid":   link.Key,
					"count": 20,
				},
			}
		} else {
			post["Component_Billboard_Billboardlist"]["next_cursor"] = nextCursor
		}
		data, _ := json.Marshal(post)

		videoList := WbVideoList{}
		page, err = w.FetchData(link, map[string]string{"data": string(data)}, headers)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(page.Content), &videoList)
		if err != nil {
			// 但凡一次报错，全部不算了
			return
		}
		if len(videoList.Data.Videos.List) == 0 {
			break
		}
		for _, v := range videoList.Data.Videos.List {
			hotList = append(hotList, Hot{
				Key:       w.FetchKey(v.Oid),
				Title:     v.Title,
				OriginUrl: fmt.Sprintf("https://weibo.com/tv/show/%s", v.Oid),
				Card:      CardVideo,
				Ext: map[string]string{
					"cover": v.Cover,
					"date":  v.Date,
					"score": v.PlayCount,
				},
			})
		}
		nextCursor = videoList.Data.Videos.Next
	}

	res = Page{
		Link: link,
		List: hotList,
		T:    time.Now(),
	}

	return
}

func (w *Wbvideo) FetchKey(key string) string {
	return key
}

func NewWbvideo() *Wbvideo {
	return &Wbvideo{
		Site{
			Name:     "微博视频",
			Key:      SITE_WBVIDEO,
			Root:     "https://weibo.com/tv/home",
			Desc:     "微博视频榜单",
			CrawType: CrawApi,
			Tabs:     WbvideoTabs,
		},
	}
}

var _ Spider = &Wbvideo{}

func init() {
	RegistSite(SITE_WBVIDEO, NewWbvideo())
}
