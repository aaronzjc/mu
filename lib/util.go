package lib

import (
	"crawler/util"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func NewSite(t string) Site {
	switch t {
	case SITE_V2EX:
		return Site{
			Name:     "v2ex",
			Key:      t,
			Root:     "https://www.v2ex.com",
			Desc:     "way to explore",
			CrawType: CrawHtml,
			Tabs:     V2exTabs,
		}
	case SITE_CT:
		return Site{
			Name:     "抽屉",
			Key:      t,
			Root:     "https://dig.chouti.com",
			Desc:     "抽屉新热榜",
			CrawType: CrawApi,
			Tabs:     ChoutiTabs,
		}
	case SITE_WEIBO:
		return Site{
			Name:     "微博",
			Key:      t,
			Root:     "https://s.weibo.com",
			Desc:     "微博热搜",
			CrawType: CrawHtml,
			Tabs:     WeiboTabs,
		}
	case SITE_ZHIHU:
		return Site{
			Name:     "知乎",
			Key:      t,
			Root:     "https://zhihu.com",
			Desc:     "知乎热榜",
			CrawType: CrawHtml,
			Tabs:     ZhihuTabs,
		}
	case SITE_HACKER:
		return Site{
			Name:     "Hacker",
			Key:      t,
			Root:     "https://news.ycombinator.com/",
			Desc:     "Hacker News",
			CrawType: CrawHtml,
			Tabs:     HackerTabs,
		}
	default:
		log.Fatalln("Unknown site name", t)
		return Site{}
	}
}

func CrawJSON(link Link) (Page, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", link.Url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[error] CrawJSON error, url = %s, err = %s\n", link.Url, err.Error())
		return Page{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var list HotList
	if err := json.Unmarshal(body, &list); err != nil {
		log.Printf(err.Error())
		return Page{}, err
	}

	return Page{
		Link:    link,
		Content: string(body),
		Json:    nil,
	}, nil
}

func CrawHTML(link Link, headers map[string]string) (Page, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", link.Url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[error] CrawHTML error, url = %s, err = %s\n", link.Url, err.Error())
		return Page{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		log.Printf("[error] Encode html error , err = %s\n", err.Error())
	}

	return Page{
		Link:    link,
		Content: bodyStr,
		Doc:     doc,
	}, nil
}

func RedisConn() *redis.Client {
	config := util.NewConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func SaveToRedis(key string, hkey string, data string) {
	client := RedisConn()
	_, err := client.HSet(key, hkey, data).Result()
	if err != nil {
		log.Printf("[error] SaveToRedis error , err = %s\n", err.Error())
	}
}
