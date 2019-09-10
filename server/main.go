package main

import (
	"crawler/lib"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

type Tag struct {
	Name 		string 			`json:"name"`
	Key 		string 			`json:"key"`
}

type Tab struct {
	Name 		string 			`json:"name"`
	Key 		string 			`json:"key"`
	Tags 		[]Tag 			`json:"tags"`
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
				Key: v["tag"],
			})
		}

		return tags
	}

	// V2ex
	tabs = append(tabs, Tab{
		Name: "v2ex",
		Key: lib.RedisV2ex,
		Tags: fetchTags(lib.V2exTabs),
	})

	// 抽屉
	tabs = append(tabs, Tab{
		Name: "抽屉",
		Key: lib.RedisCt,
		Tags: fetchTags(lib.ChoutiTabs),
	})

	// 知乎
	tabs = append(tabs, Tab{
		Name: "知乎",
		Key: lib.RedisZhihu,
		Tags: fetchTags(lib.ZhihuTabs),
	})

	// 微博
	tabs = append(tabs, Tab{
		Name: "微博",
		Key: lib.RedisWeibo,
		Tags: fetchTags(lib.WeiboTabs),
	})

	data, _ := json.Marshal(tabs)

	JSON(w, data)
}

func aj(w http.ResponseWriter, req *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr: "10.8.77.119:6379",
		Password: "",
		DB: 0,
	})
	key := req.URL.Query()["key"][0]
	hkey := req.URL.Query()["hkey"][0]

	data, err := client.HGet(key, hkey).Result()

	if err != nil {
		log.Println("[info] aj req empty " + err.Error())
		JSON(w, []byte(`{"list": []}`))
		return
	}

	var hotJson lib.HotJson
	err = json.Unmarshal([]byte(data), &hotJson)
	if err != nil {
		log.Println("[error] aj req error " + err.Error())
		JSON(w, []byte(`{"list": []}`))
		return
	}

	js, _ := json.Marshal(hotJson)
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
