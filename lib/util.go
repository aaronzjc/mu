package lib

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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
		Link: link,
		Content: string(body),
		Json: nil,
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
		Link: link,
		Content: bodyStr,
		Doc: doc,
	}, nil
}

func RedisConn() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "10.8.77.119:6379",
		Password: "",
		DB: 0,
	})
}

func SaveToRedis(key string, hkey string, data string) {
	client := RedisConn()
	client.HSet(key, hkey, data)
}