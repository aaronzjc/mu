package lib

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"mu/internal/util/logger"
	"net/http"
	"strings"
	"time"
)

const (
	CrawApi = iota + 1
	CrawHtml
)

// 热榜新闻
type Hot struct {
	Id        int     `json:"id"`
	Key       string  `json:"key"`
	Title     string  `json:"title"`
	Rank      float64 `json:"rank"`
	OriginUrl string  `json:"origin_url"`
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
	Key string
	Url string
	Tag string
	Sp  Spider
}

// 抓取的页面信息
type Page struct {
	Link    Link
	Content string

	Doc  *goquery.Document
	Json []map[string]interface{}

	List []Hot
	T    time.Time
}

func CrawJSON(link Link) (Page, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", link.Url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("CrawJSON error, url = %s, err = %s\n", link.Url, err.Error())
		return Page{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var list HotList
	if err := json.Unmarshal(body, &list); err != nil {
		logger.Error(err.Error())
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
	if err != nil {
		logger.Error("CrawHTML error, url = %s, err = %s\n", link.Url, err.Error())
		return Page{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		logger.Error("Encode html error , err = %s\n", err.Error())
	}

	return Page{
		Link:    link,
		Content: bodyStr,
		Doc:     doc,
	}, nil
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

func (s *Site) Craw(link Link, headers map[string]string) (Page, error) {
	var page Page
	var err error
	if s.CrawType == CrawApi {
		page, err = CrawJSON(link)
	} else if s.CrawType == CrawHtml {
		page, err = CrawHTML(link, headers)
	} else {
		err = errors.New("[error] No matched CrawType")
	}
	if err != nil {
		return Page{}, err
	}

	return page, nil
}

func FSite(t string) Spider {
	switch t {
	case SITE_V2EX:
		return &V2ex{NewSite(t)}
	case SITE_CT:
		return &Chouti{NewSite(t)}
	case SITE_WEIBO:
		return &Weibo{NewSite(t)}
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
			Root:     "http://www.guanggoo.com",
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
		SITE_GUANGGU,
		SITE_HACKER,
		SITE_GITHUB,
		SITE_TIEBA,
		SITE_REDDIT,
		SITE_ZAOBAO,
	}
}
