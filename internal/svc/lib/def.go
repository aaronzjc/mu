package lib

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"mu/internal/util/logger"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	CrawApi = iota + 1
	CrawHtml
)

const (
	CardText = iota
	CardRichText
	CardVideo
)

// 热榜
type Hot struct {
	Key       string            `json:"key"`
	Title     string            `json:"title"`
	Desc      string            `json:"desc"`
	Rank      float64           `json:"rank"`
	OriginUrl string            `json:"origin_url"`
	Card      uint8             `json:"card_type"`
	Ext       map[string]string `json:"ext"`
}

// 热榜新闻列表
type HotJson struct {
	T    string `json:"t"`
	List []Hot  `json:"list"`
}

type Spider interface {
	BuildUrl() ([]Link, error)
	CrawPage(link Link, headers map[string]string) (Page, error)
	FetchKey(link string) string
}

// 链接信息
type Link struct {
	Key    string
	Url    string
	Tag    string
	Method string

	Sp Spider
}

// 抓取的页面信息
type Page struct {
	Link Link

	Content string
	Doc     *goquery.Document
	Json    []map[string]interface{}

	List []Hot
	T    time.Time
}

// 站点信息
type Site struct {
	Name     string
	Key      string
	Root     string
	Desc     string
	CrawType int8
	Tabs     []map[string]string
}

func (s *Site) FetchData(link Link, params map[string]string, headers map[string]string) (res Page, err error) {
	var data io.Reader
	// 构造请求参数
	switch link.Method {
	case "GET":
	case "POST":
		var contentType string
		contentType, ok := headers["Content-Type"]
		if !ok {
			contentType = "application/json"
			if headers == nil {
				headers = make(map[string]string)
			}
			headers["Content-Type"] = "application/json"
		}
		switch contentType {
		case "application/json":
			var jstr []byte
			jstr, err = json.Marshal(params)
			data = bytes.NewReader(jstr)
		case "application/x-www-form-urlencoded":
			dstr := url.Values{}
			for k, v := range params {
				dstr.Add(k, v)
			}
			data = strings.NewReader(dstr.Encode())
		}
	default:
		link.Method = "GET"
	}

	if err != nil {
		return
	}

	client := http.Client{}
	req, err := http.NewRequest(link.Method, link.Url, data)
	if err != nil {
		logger.Error("init http request error e = %v", err)
		return
	}

	// 添加浏览器头
	if _, ok := headers["User-Agent"]; !ok {
		req.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	}
	if len(headers) > 0 {
		for k, v := range headers {
			if k == "" || v == "" {
				continue
			}
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logger.Error("request error e = %v", err)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	var doc *goquery.Document
	if s.CrawType == CrawHtml {
		doc, err = goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
		if err != nil {
			return
		}
	}

	return Page{
		Link:    link,
		Content: bodyStr,
		Doc:     doc,
	}, nil
}

func FSite(t string) Spider {
	switch t {
	case SITE_V2EX:
		return &V2ex{NewSite(t)}
	case SITE_CT:
		return &Chouti{NewSite(t)}
	case SITE_WEIBO:
		return &Weibo{NewSite(t)}
	case SITE_WBVIDEO:
		return &Wbvideo{NewSite(t)}
	case SITE_ZHIHU:
		return &Zhihu{NewSite(t)}
	case SITE_GUANGGU:
		return &Guanggu{NewSite(t)}
	case SITE_HACKER:
		return &Hacker{NewSite(t)}
	case SITE_GITHUB:
		return &Github{NewSite(t)}
	case SITE_TIEBA:
		return &Tieba{NewSite(t)}
	case SITE_REDDIT:
		return &Reddit{NewSite(t)}
	case SITE_ZAOBAO:
		return &Zaobao{NewSite(t)}
	default:
		logger.Fatal("Unknown site name " + t)
		return nil
	}
}

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
	case SITE_WBVIDEO:
		return Site{
			Name:     "微博视频",
			Key:      t,
			Root:     "https://weibo.com/tv/home",
			Desc:     "微博视频榜单",
			CrawType: CrawApi,
			Tabs:     WbvideoTabs,
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
	case SITE_GUANGGU:
		return Site{
			Name:     "光谷",
			Key:      t,
			Root:     "https://www.guozaoke.com",
			Desc:     "武汉光谷社区",
			CrawType: CrawHtml,
			Tabs:     GuangGuTabs,
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
	case SITE_GITHUB:
		return Site{
			Name:     "Github",
			Key:      t,
			Root:     "https://github.com",
			Desc:     "Github.com",
			CrawType: CrawHtml,
			Tabs:     GithubTabs,
		}
	case SITE_TIEBA:
		return Site{
			Name:     "贴吧",
			Key:      t,
			Root:     "https://tieba.baidu.com",
			Desc:     "鱼龙混杂的社区",
			CrawType: CrawHtml,
			Tabs:     TiebaTabs,
		}
	case SITE_REDDIT:
		return Site{
			Name:     "Reddit",
			Key:      t,
			Root:     "https://www.reddit.com/",
			Desc:     "老外的贴吧",
			CrawType: CrawApi,
			Tabs:     RedditTabs,
		}
	case SITE_ZAOBAO:
		return Site{
			Name:     "联合早报",
			Key:      t,
			Root:     "http://www.zaobao.com",
			Desc:     "新加坡新闻",
			CrawType: CrawHtml,
			Tabs:     ZaobaoTabs,
		}
	default:
		logger.Fatal("Unknown site name " + t)
		return Site{}
	}
}

func AvailableSites() []string {
	return []string{
		SITE_V2EX,
		SITE_CT,
		SITE_ZHIHU,
		SITE_WEIBO,
		SITE_WBVIDEO,
		SITE_GUANGGU,
		SITE_HACKER,
		SITE_GITHUB,
		SITE_TIEBA,
		SITE_REDDIT,
		SITE_ZAOBAO,
	}
}
