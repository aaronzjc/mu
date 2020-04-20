package model

import (
	"encoding/json"
	"errors"
	"mu/internal/svc/lib"
	"strings"
)

type CrawType int
type NodeOption int8
type Status int8

const (
	CrawHtml CrawType = 1 // 网站是HTML
	CrawApi  CrawType = 2 // 网站是JSON接口

	ByType  NodeOption = 1 // 通过服务器类型
	ByHosts NodeOption = 2 // 服务器IPs

	Disable Status = 0 // 禁用
	Enable  Status = 1 // 启用
)

type Site struct {
	ID         int
	Name       string     `gorm:"name"`
	Root       string     `gorm:"root"`
	Key        string     `gorm:"key"`
	Desc       string     `gorm:"desc"`
	Type       int8       `gorm:"type"`
	Tags       string     `gorm:"tags"`
	Cron       string     `gorm:"cron"`
	Enable     Status     `gorm:"enable"`
	NodeOption NodeOption `gorm:"node_option"`
	NodeType   int8       `gorm:"node_type"`
	NodeHosts  string     `gorm:"node_hosts"`
	ReqHeaders string     `gorm:"req_headers"`
}

type Tag struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Enable int8   `json:"enable"`
}

type Header struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type SiteJson struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Key        string     `json:"key"`
	Root       string     `json:"root"`
	Desc       string     `json:"desc"`
	Tags       []Tag      `json:"tags"`
	Type       int8       `json:"type"`
	Cron       string     `json:"cron"`
	NodeOption NodeOption `json:"node_option"`
	NodeType   int8       `json:"node_type"`
	NodeHosts  []int      `json:"node_hosts"`
	Enable     Status     `json:"enable"`
	ReqHeaders []Header   `json:"req_headers"`
}

func (s *Site) TableName() string {
	return "site"
}

func (s *Site) CheckArgs() error {
	if s.Name == "" {
		return errors.New("名字为空")
	}
	if len(strings.Split(s.Cron, " ")) != 5 {
		return errors.New("cron必须是5位表达式")
	}

	return nil
}

func (s *Site) Create() error {
	var tmp []Site
	err := FetchRows(Query{
		Query: "`key` = ? or `root` = ?",
		Args:  []interface{}{s.Key, s.Root},
	}, &tmp)
	if err != nil {
		return errors.New("create site error")
	}
	if len(tmp) > 0 {
		return errors.New("同key或者同root的站点已经存在")
	}

	err = Create(&s)
	if err != nil {
		return errors.New("create site err")
	}

	return nil
}

func (s *Site) Update(data map[string]interface{}) error {
	err := Update(&s, data)
	if err != nil {
		return errors.New("update site failed")
	}

	return nil
}

func (s *Site) FetchInfo() (Site, error) {
	var tmp Site

	err := First(Query{
		Query: "`id` = ?",
		Args:  []interface{}{s.ID},
	}, &tmp)
	if err != nil {
		return Site{}, errors.New("fetch site info failed")
	}

	return tmp, nil
}

func (s *Site) FetchRow(query Query) (Site, error) {
	var site Site
	err := First(query, &site)
	if err != nil {
		return Site{}, errors.New("fetchRow site failed")
	}
	return site, nil
}

func (s *Site) FetchRows(query Query) ([]Site, error) {
	var list []Site
	err := FetchRows(query, &list)
	if err != nil {
		return nil, errors.New("fetchRows site failed")
	}
	return list, nil
}

func (s *Site) FormatJson() (SiteJson, error) {
	tags := []Tag{}
	headers := []Header{}
	hosts := []int{}

	var err error
	if s.Tags != "" {
		err = json.Unmarshal([]byte(s.Tags), &tags)
		if err != nil {
			return SiteJson{}, errors.New("标签解析失败")
		}
	}

	if s.ReqHeaders != "" {
		err = json.Unmarshal([]byte(s.ReqHeaders), &headers)
		if err != nil {
			return SiteJson{}, errors.New("头解析失败")
		}
	}

	if s.NodeHosts != "" {
		err = json.Unmarshal([]byte(s.NodeHosts), &hosts)
		if err != nil {
			return SiteJson{}, errors.New("节点解析失败")
		}
	}

	return SiteJson{
		ID:         s.ID,
		Name:       s.Name,
		Key:        s.Key,
		Root:       s.Root,
		Desc:       s.Desc,
		Tags:       tags,
		Type:       s.Type,
		Cron:       s.Cron,
		NodeOption: s.NodeOption,
		NodeType:   s.NodeType,
		NodeHosts:  hosts,
		Enable:     s.Enable,
		ReqHeaders: headers,
	}, nil
}

func (s *Site) InitSites() {
	var tagStr []byte

	avaSites := lib.AvailableSites()
	for _, siteKey := range avaSites {
		site := lib.NewSite(siteKey)
		row, err := s.FetchRow(Query{
			Query: " `key` = ? ",
			Args:  []interface{}{site.Key},
		})
		if err != nil {
			panic("init sites fetch failed " + err.Error())
		}

		tagMap := make(map[string]int)
		rowJson, _ := row.FormatJson()
		for _, v := range rowJson.Tags {
			if v.Enable == 0 {
				tagMap[v.Key] = 1
			}
		}
		var tags []Tag
		var e int8
		for _, tag := range site.Tabs {
			e = 1
			if _, ok := tagMap[tag["tag"]]; ok {
				e = 0
			}
			tags = append(tags, Tag{
				Key:    tag["tag"],
				Name:   tag["name"],
				Enable: e,
			})
		}
		tagStr, _ = json.Marshal(tags)

		if row.ID > 0 {
			err = row.Update(map[string]interface{}{
				"name": site.Name,
				"root": site.Root,
				"tags": string(tagStr),
				"type": site.CrawType,
			})
			if err != nil {
				panic("init sites update failed " + err.Error())
			}
			continue
		}

		row = Site{
			Name:       site.Name,
			Key:        site.Key,
			Root:       site.Root,
			Cron:       "*/30 * * * *",
			NodeOption: 1, // 默认使用分类
			NodeType:   1, // 默认国内的机器
			NodeHosts:  "",
			Desc:       site.Desc,
			Tags:       string(tagStr),
			Type:       site.CrawType,
			ReqHeaders: "",
		}
		err = row.Create()
		if err != nil {
			panic("init sites create failed " + err.Error())
		}
	}
}

func (s *Site) FixNodeId(delId int) {
	sites, _ := (&Site{}).FetchRows(Query{})
	for _, site := range sites {
		sj, _ := site.FormatJson()
		var newHosts []int
		for _, v := range sj.NodeHosts {
			if v == delId {
				continue
			}
			newHosts = append(newHosts, v)
		}
		jstr, _ := json.Marshal(newHosts)
		_ = site.Update(map[string]interface{}{
			"node_hosts": jstr,
		})
	}
}
