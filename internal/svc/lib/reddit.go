package lib

import (
	"encoding/json"
	"fmt"
	"mu/internal/util/logger"
	"regexp"
	"time"
)

const SITE_REDDIT = "reddit"

var RedditTabs = []map[string]string{
	{
		"tag":  "AskReddit",
		"name": "AskReddit",
		"url":  "AskReddit",
	},
	{
		"tag":  "Jokes",
		"name": "Jokes",
		"url":  "Jokes",
	},
	{
		"tag":  "leagueoflegends",
		"name": "lol",
		"url":  "leagueoflegends",
	},
}

type Reddit struct {
	Site
}

type RedditList struct {
	Posts map[string]map[string]interface{} `json:"posts"`
}

func (r *Reddit) BuildUrl() ([]Link, error) {
	str := "https://gateway.reddit.com/desktopapi/v1/subreddits/%s?rtj=only&redditWebClient=web2x&app=web2x-client-production&allow_over18=&include=prefsSubreddit&after=&dist=11&layout=card&sort=hot"

	var list []Link
	for _, tab := range RedditTabs {
		url := tab["url"]
		link := Link{
			Key: url,
			Url: fmt.Sprintf(str, url),
			Tag: tab["tag"],
			Sp:  r,
		}
		list = append(list, link)
	}

	return list, nil
}

func (r *Reddit) CrawPage(link Link) (Page, error) {
	page, err := r.Craw(link, nil)
	if err != nil {
		return Page{}, err
	}

	var list RedditList
	if err := json.Unmarshal([]byte(page.Content), &list); err != nil {
		logger.Error("%v", err.Error())
		return Page{}, err
	}

	var listArr []map[string]interface{}
	for k, v := range list.Posts {
		if ok, _ := regexp.Match("redditads", []byte(v["permalink"].(string))); ok {
			continue
		}
		listArr = append(listArr, map[string]interface{}{
			"key":   k,
			"title": v["title"],
			"url":   v["permalink"],
			"score": v["score"],
		})
	}

	page.Json = listArr

	var data []Hot
	for _, v := range page.Json {
		h := Hot{
			Id:        0,
			Title:     v["title"].(string),
			OriginUrl: v["url"].(string),
			Rank:      v["score"].(float64),
			Key:       v["key"].(string),
		}
		if h.Key == "" {
			continue
		}
		data = append(data, h)
	}

	page.T = time.Now()
	page.List = data

	return page, nil
}

func (r *Reddit) FetchKey(key string) string {
	return key
}
