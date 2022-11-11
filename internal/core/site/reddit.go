package site

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/aaronzjc/mu/pkg/logger"
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

type RedditList struct {
	Posts map[string]map[string]interface{} `json:"posts"`
}

type Reddit struct {
	Site
}

var _ Spider = &Reddit{}

func init() {
	RegistSite(SITE_REDDIT, &Reddit{
		Site{
			Name:     "Reddit",
			Key:      SITE_REDDIT,
			Root:     "https://www.reddit.com/",
			Desc:     "老外的贴吧",
			CrawType: CrawApi,
			Tabs:     RedditTabs,
		},
	})
}

func (r *Reddit) GetSite() *Site {
	return &r.Site
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
		}
		list = append(list, link)
	}

	return list, nil
}

func (r *Reddit) CrawPage(link Link, headers map[string]string) (Page, error) {
	page, err := r.FetchData(link, nil, nil)
	if err != nil {
		return Page{}, err
	}

	var list RedditList
	if err := json.Unmarshal([]byte(page.Content), &list); err != nil {
		logger.Error("%v", err.Error())
		return Page{}, err
	}

	var listArr []map[string]interface{}
	re, _ := regexp.Compile(`redditads`)
	for k, v := range list.Posts {
		if ok := re.Match([]byte(v["permalink"].(string))); ok {
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
