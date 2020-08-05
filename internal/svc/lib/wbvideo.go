package lib

import (
	"encoding/json"
	"fmt"
	"mu/internal/util/tool"
	"time"
)

const SITE_WBVIDEO = "wbvideo"

var WbvideoTabs = []map[string]string{
	{
		"tag":  "all",
		"url":  "https://weibo.com/tv/api/component?page=%2Ftv%2Fbillboard%2F4418219809678881",
		"name": "全站",
	},
}

type Wbvideo struct {
	Site
}

type WbVideoList struct {
	Code string `json:"code"`
	Msg	 string `json:"msg"`
	Data struct{
		Categorys interface{} `json:"Component_Billboard_Billboardcategory"`
		Videos struct{
			Next int `json:"next_cursor"`
			List []struct{
				Title string `json:"title"`
				Cover string `json:"cover_image"`
				Id 	 string `json:"mid"`
				Oid string `json:"oid"`
				Date string `json:"date"`
				PlayCount string `json:"play_count"`
			} `json:"list"`
		} `json:"Component_Billboard_Billboardlist"`
	} `json:"data"`
}

func (w *Wbvideo) BuildUrl() ([]Link, error) {
	var list []Link
	for _, tab := range WeiboTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: url,
			Tag: tab["tag"],
			Sp:  w,
		}
		list = append(list, link)
	}

	return list, nil
}

func (w *Wbvideo) CrawPage(link Link, headers map[string]string) (res Page, err error) {
	var videos []map[string]interface{}
	var nextCursor string
	post := make(map[string]map[string]interface{})
	for {
		if nextCursor == "" {
			post = map[string]map[string]interface{}{
				"Component_Billboard_Billboardcategory": {},
				"Component_Billboard_Billboardlist": {
					"cid": "4418213501411061",
					"count": 20,
				},
			}
		} else {
			post["Component_Billboard_Billboardlist"]["next_cursor"] = nextCursor
		}
		data, _ := json.Marshal(post)

		videoList := WbVideoList{}
		page, err := w.FetchData(link, map[string]string{"data": string(data)}, headers)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(page.Content), videoList)
		if err != nil {
			// 但凡一次报错，全部不算了
			return
		}
		if len(videoList.Data.Videos.List) == 0 {
			break
		}
		for _, v := range videoList.Data.Videos.List {
			videos = append(videos, map[string]interface{}{
				"id": v.Id,
				"url": fmt.Sprintf("https://weibo.com/tv/show/%s", v.Oid),
				"title": v.Title,
				"cover": v.Cover,
				"date": v.Date,
				"count": v.PlayCount,
			})
		}
	}
	
	res = Page{
		Link:    link,
		Content: "",
		Doc:     nil,
		Json:    nil,
		List:    nil,
		T:       time.Time{},
	}

	return
}

func (w *Wbvideo) FetchKey(link string) string {
	if link == "" {
		return ""
	}
	return tool.MD55(link)
}
