package lib

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"time"
)

const (
	CrawApi = iota + 1
	CrawHtml
)

type Site struct {
	Name     string
	Key      string
	Root     string
	Desc     string
	CrawType int
	Tabs     []map[string]string
}

type Link struct {
	Key string
	Url string
	Tag string
	Sp  Spider
}

type Page struct {
	Link    Link
	Content string

	Doc  *goquery.Document
	Json []map[string]interface{}

	List []Hot
	T    time.Time
}

type Hot struct {
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Rank      float64 `json:"rank"`
	OriginUrl string  `json:"origin_url"`
}

type HotJson struct {
	T    string `json:"t"`
	List []Hot  `json:"list"`
}

type Spider interface {
	BuildUrl() ([]Link, error)
	CrawPage(link Link) (Page, error)
	Store(page Page) bool
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
