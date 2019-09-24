package model

import (
	"crawler/internal/lib"
	"encoding/json"
	"errors"
	"log"
	"strings"
)

type CrawType int

const (
	_        CrawType = iota
	CrawHtml          // 1
	CrawApi           // 2
)

type Site struct {
	ID        int
	Name      string    `gorm:"name"`
	Root      string    `gorm:"root"`
	Key       string    `gorm:"key"`
	Desc      string    `gorm:"desc"`
	Type      int8      `gorm:"type"`
	Tags      string    `gorm:"tags"`
	Cron      string    `gorm:"cron"`
	Enable    int8      `gorm:"enable"`
	NodeOption int8		`gorm:"node_option"`
	NodeType  int8    `gorm:"node_type"`
	NodeHosts string    `gorm:"node_hosts"`
}

type Tag struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Enable int8   `json:"enable"`
}

type SiteJson struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Key        string   `json:"key"`
	Root       string   `json:"root"`
	Desc       string   `json:"desc"`
	Tags       []Tag    `json:"tags"`
	Type       int8 	`json:"type"`
	Cron       string   `json:"cron"`
	NodeOption int8     `json:"node_option"`
	NodeType   int8     `json:"node_type"`
	NodeHosts  []int    `json:"node_hosts"`
	Enable     int8     `json:"enable"`
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
	tmp, err := s.FetchRows("`key` = ? or `root` = ?", s.Key, s.Root)
	if err != nil {
		return errors.New("create site error " + err.Error())
	}
	if len(tmp) > 0 {
		return errors.New("同key或者同root的站点已经存在")
	}

	db := DPool().Conn
	db = db.Create(&s)
	if err = db.Error; err != nil {
		log.Printf("[error] create err %v, exp %s \n", err, db.QueryExpr())
		return errors.New("create site err")
	}

	return nil
}

func (s *Site) Update(data map[string]interface{}) error {
	db := DPool().Conn

	db = db.Model(&s).Update(data)
	if err := db.Error; err != nil {
		log.Printf("[error] update err %v, exp %s \n", err, db.QueryExpr())
		return errors.New("update site failed")
	}

	return nil
}

func (s *Site) FetchInfo() (Site, error) {
	var tmp Site
	db := DPool().Conn
	db = db.Where("id = ?", s.ID).First(&tmp)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		log.Printf("[error] FetchInfo err %v, exp %s \n", err, db.QueryExpr())
		return Site{}, errors.New("fetch site info failed")
	}

	return tmp, nil
}

func (s *Site) FetchRow(query string, args ...interface{}) (Site, error) {
	db := DPool().Conn

	var site Site
	db = db.Where(query, args...).First(&site)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		log.Printf("[error] FetchRows err %v, exp %s \n", err, db.QueryExpr())
		return Site{}, errors.New("fetchRow site failed")
	}
	return site, nil
}

func (s *Site) FetchRows(query string, args ...interface{}) ([]Site, error) {
	db := DPool().Conn

	var list []Site
	db = db.Where(query, args...).Find(&list)
	if err := db.Error; err != nil {
		log.Printf("[error] FetchRows err %v, exp %s \n", err, db.QueryExpr())
		return nil, errors.New("fetchRows site failed")
	}
	return list, nil
}

func (s *Site) FormatJson() (SiteJson, error) {
	var tags []Tag
	var hosts []int

	var err error
	if s.Tags != "" {
		err = json.Unmarshal([]byte(s.Tags), &tags)
		if err != nil {
			return SiteJson{}, errors.New("标签解析失败")
		}
	} else {
		tags = []Tag{}
	}

	if s.NodeHosts != "" {
		err = json.Unmarshal([]byte(s.NodeHosts), &hosts)
		if err != nil {
			return SiteJson{}, errors.New("节点解析失败")
		}
	} else {
		hosts = []int{}
	}

	return SiteJson{
		ID: s.ID,
		Name: s.Name,
		Key: s.Key,
		Root: s.Root,
		Desc: s.Desc,
		Tags: tags,
		Type: s.Type,
		Cron: s.Cron,
		NodeOption: s.NodeOption,
		NodeType: s.NodeType,
		NodeHosts: hosts,
		Enable: s.Enable,
	}, nil
}

func (s *Site) InitSites() {
	var tagStr []byte

	avaSites := lib.AvailableSites()
	for _, siteKey := range avaSites {
		site := lib.NewSite(siteKey)
		row, err := s.FetchRow(" `key` = ? ", site.Key)
		if err != nil {
			panic("init sites fetch failed " + err.Error())
		}

		var tags []Tag
		for _, tag := range site.Tabs {
			tags = append(tags, Tag{
				Key: tag["tag"],
				Name: tag["name"],
				Enable: 1,
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
			Name: site.Name,
			Key: site.Key,
			Root: site.Root,
			Desc: site.Desc,
			Tags: string(tagStr),
			Type: site.CrawType,
		}
		err = row.Create()
		if err != nil {
			panic("init sites create failed " + err.Error())
		}
	}
}