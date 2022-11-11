package site

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aaronzjc/mu/pkg/logger"
)

const (
	// 抓取类型
	CrawApi = iota + 1
	CrawHtml

	// 渲染的样式
	CardText = iota
	CardRichText
	CardVideo
)

var (
	// 网站配置
	SiteMap = make(map[string]Spider)
	// 可用网站
	Avaiable = []string{
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

type Spider interface {
	GetSite() *Site
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
	if err != nil {
		logger.Error("request error e = %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
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

func RegistSite(name string, s Spider) {
	SiteMap[name] = s
}
