package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type NodeType int8
type Ping int8

const (
	TypeOverseas NodeType = 1 // 海外
	TypeMainland NodeType = 2 // 大陆

	PingFailed Ping = 0
	PingOk     Ping = 1
)

type Node struct {
	ID       int       `gorm:"id" json:"id"`
	Name     string    `gorm:"name" json:"name"`
	Addr     string    `gorm:"addr" json:"addr"`
	Type     int8      `gorm:"type" json:"type"`
	Enable   int8      `gorm:"enable" json:"enable"`
	Ping     Ping      `gorm:"ping" json:"ping"`
	CreateAt time.Time `gorm:"create_at" json:"create_at"`
}

func (node *Node) TableName() string {
	return "node"
}

func (node *Node) CheckArgs() error {
	if node.Name == "" {
		return errors.New("参数为空")
	}
	if match, _ := regexp.MatchString("^(\\d+)\\.(\\d+)\\.(\\d+)\\.(\\d+):\\d+$", node.Addr); !match {
		return errors.New("Addr不规范")
	}

	return nil
}

func (node *Node) Create() error {
	tmp, err := node.FetchInfo()
	if err != nil {
		return errors.New("create node err")
	}

	if tmp.ID > 0 {
		return errors.New(fmt.Sprintf("node with %s exists", node.Addr))
	}

	err = Create(&node)
	if err != nil {
		return errors.New("create node err")
	}

	return nil
}

func (node *Node) Del() bool {
	return Del(&node) == nil
}

func (node *Node) Update(data map[string]interface{}) error {
	err := Update(&node, data)
	if err != nil {
		return errors.New("update node failed")
	}

	return nil
}

func (node *Node) FetchInfo() (Node, error) {
	var n Node

	err := First(Query{
		Query: "`id` = ?",
		Args:  []interface{}{node.ID},
	}, &n)
	if err != nil {
		return Node{}, errors.New("fetch node info failed")
	}

	return n, nil
}

func (node *Node) FetchRows(query Query) ([]Node, error) {
	var list []Node

	err := FetchRows(query, &list)
	if err != nil {
		return nil, errors.New("fetchrows node failed")
	}

	return list, nil
}

func (node *Node) FormatJson() (Node, error) {
	json := Node{
		ID:       node.ID,
		Name:     node.Name,
		Addr:     node.Addr,
		Type:     node.Type,
		Enable:   node.Enable,
		Ping:     node.Ping,
		CreateAt: time.Now(),
	}

	return json, nil
}
