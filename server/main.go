package main

import (
	"crawler/lib"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"time"
)

type Tag struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Tab struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Tags []Tag  `json:"tags"`
}

func JSON(w http.ResponseWriter, data []byte) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

func config(w http.ResponseWriter, req *http.Request) {
	var tabs []Tab

	var fetchTags = func(tabs []map[string]string) []Tag {
		var tags []Tag
		for _, v := range tabs {
			tags = append(tags, Tag{
				Name: v["name"],
				Key:  v["tag"],
			})
		}

		return tags
	}

	// V2ex
	tabs = append(tabs, Tab{
		Name: "v2ex",
		Key:  lib.RedisV2ex,
		Tags: fetchTags(lib.V2exTabs),
	})

	// 抽屉
	tabs = append(tabs, Tab{
		Name: "抽屉",
		Key:  lib.RedisCt,
		Tags: fetchTags(lib.ChoutiTabs),
	})

	// 知乎
	tabs = append(tabs, Tab{
		Name: "知乎",
		Key:  lib.RedisZhihu,
		Tags: fetchTags(lib.ZhihuTabs),
	})

	// 微博
	tabs = append(tabs, Tab{
		Name: "微博",
		Key:  lib.RedisWeibo,
		Tags: fetchTags(lib.WeiboTabs),
	})

	// HackerNews
	tabs = append(tabs, Tab{
		Name: "Hacker",
		Key:  lib.RedisHacker,
		Tags: fetchTags(lib.HackerTabs),
	})

	data, _ := json.Marshal(tabs)

	JSON(w, data)
}

func aj(w http.ResponseWriter, req *http.Request) {
	t := time.Now()

	client := redis.NewClient(&redis.Options{
		Addr:     "10.8.77.119:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	elapsed := time.Since(t)
	log.Println("redis connect", elapsed)

	key := req.URL.Query()["key"][0]
	hkey := req.URL.Query()["hkey"][0]

	data, err := client.HGet(key, hkey).Result()

	elapsed = time.Since(t)
	log.Println("redis hget", elapsed)

	if err != nil {
		log.Println("[info] aj req empty " + err.Error())
		JSON(w, []byte(`{"list": [], "t":""}`))
		return
	}

	var hotJson lib.HotJson
	err = json.Unmarshal([]byte(data), &hotJson)
	if err != nil {
		log.Println("[error] aj req error " + err.Error())
		JSON(w, []byte(`{"list": [], "t":""}`))
		return
	}

	js, _ := json.Marshal(hotJson)
	elapsed = time.Since(t)
	log.Println("Marshal data", elapsed)
	JSON(w, []byte(js))
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/aj", aj)
	http.HandleFunc("/config", config)

	log.Println("listen on :7980")

	log.Fatal(http.ListenAndServe(":7980", nil))
}
