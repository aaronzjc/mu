package dto

import (
	"encoding/json"

	"github.com/aaronzjc/mu/internal/domain/model"
)

type Tag struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Enable int8   `json:"enable"`
}

type Header struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

const (
	CrawHtml int8 = 1 // 网站是HTML
	CrawApi  int8 = 2 // 网站是JSON接口

	ByType  int8 = 1 // 通过服务器类型
	ByHosts int8 = 2 // 服务器IPs

	Disable int8 = 0 // 禁用
	Enable  int8 = 1 // 启用
)

type Site struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Key        string   `json:"key"`
	Root       string   `json:"root"`
	Desc       string   `json:"desc"`
	Tags       []Tag    `json:"tags"`
	Type       int8     `json:"type"`
	Cron       string   `json:"cron"`
	NodeOption int8     `json:"node_option"`
	NodeType   int8     `json:"node_type"`
	NodeHosts  []int    `json:"node_hosts"`
	Enable     int8     `json:"enable"`
	ReqHeaders []Header `json:"req_headers"`
}

func (s *Site) FillByModel(site model.Site) *Site {
	s.ID = site.ID
	s.Name = site.Name
	s.Key = site.Key
	s.Root = site.Root
	s.Desc = site.Desc
	s.Type = site.Type
	s.Cron = site.Cron
	s.NodeOption = site.NodeOption
	s.NodeType = site.NodeType
	s.Enable = site.Enable

	var err error
	tags := []Tag{}
	headers := []Header{}
	hosts := []int{}
	s.Tags = []Tag{}
	if site.Tags != "" {
		if err = json.Unmarshal([]byte(site.Tags), &tags); err == nil {
			s.Tags = tags
		}
	}
	s.ReqHeaders = []Header{}
	if site.ReqHeaders != "" {
		if err = json.Unmarshal([]byte(site.ReqHeaders), &headers); err == nil {
			s.ReqHeaders = headers
		}
	}
	s.NodeHosts = []int{}
	if site.NodeHosts != "" {
		if err = json.Unmarshal([]byte(site.NodeHosts), &hosts); err == nil {
			s.NodeHosts = hosts
		}
	}

	return s
}

type IndexSite struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Tags []Tag  `json:"tags"`
}

type NewsItem struct {
	Key       string            `json:"key"`
	Title     string            `json:"title"`
	Desc      string            `json:"desc"`
	Rank      float64           `json:"rank"`
	OriginUrl string            `json:"origin_url"`
	Card      uint8             `json:"card_type"`
	Ext       map[string]string `json:"ext"`
	Mark      bool              `json:"mark"`
}
type News struct {
	T    string     `json:"t"`
	List []NewsItem `json:"list"`
}
