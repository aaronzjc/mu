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

func config(w http.ResponseWriter, req *http.Request) {
	var tabs []Tab
	// V2ex
	v2ex := Tab{
		Name: "v2ex",
		Key: lib.RedisV2ex,
		Tags: (func() []Tag {
			var tags []Tag
			for _, v := range lib.V2exTabs {
				tags = append(tags, Tag{
					Name: v["name"],
					Key: v["tag"],
				})
			}
			return tags
		})(),
	}
	tabs = append(tabs, v2ex)

	// 抽屉
	chouti := Tab{
		Name: "抽屉",
		Key: lib.RedisCt,
		Tags: (func() []Tag {
			var tags []Tag
			for _, v := range lib.ChoutiTabs {
				tags = append(tags, Tag{
					Name: v["name"],
					Key: v["tag"],
				})
			}
			return tags
		})(),
	}
	tabs = append(tabs, chouti)

	// 知乎
	zhihu := Tab{
		Name: "知乎",
		Key: lib.RedisZhihu,
		Tags: (func() []Tag {
			var tags []Tag
			for _, v := range lib.ZhihuTabs {
				tags = append(tags, Tag{
					Name: v["name"],
					Key: v["tag"],
				})
			}
			return tags
		})(),
	}
	tabs = append(tabs, zhihu)

	// 微博
	weibo := Tab{
		Name: "抽屉",
		Key: lib.RedisWeibo,
		Tags: (func() []Tag {
			var tags []Tag
			for _, v := range lib.WeiboTabs {
				tags = append(tags, Tag{
					Name: v["name"],
					Key: v["tag"],
				})
			}
			return tags
		})(),
	}
	tabs = append(tabs, weibo)

	data, _ := json.Marshal(tabs)

	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(data))
}

func aj(w http.ResponseWriter, req *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr: "10.8.77.119:6379",
	})
	key := req.URL.Query()["key"][0]
	hkey := req.URL.Query()["hkey"][0]

	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := client.HGet(key, hkey).Result()

	if err != nil {
		log.Fatalf("[info] aj req empty ")
		w.Write([]byte(`{"list": []}`))
		return
	}

	var hotJson lib.HotJson
	err = json.Unmarshal([]byte(data), &hotJson)
	if err != nil {
		log.Fatalf("[error] aj req error " + err.Error())
		w.Write([]byte(`{"list": []}`))
		return
	}

	js, _ := json.Marshal(hotJson)
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(js))
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/aj", aj)
	http.HandleFunc("/config", config)

	log.Println("listen on :7980")

	log.Fatal(http.ListenAndServe(":7980", nil))
}
