package model

import (
	"crawler/internal/util/logger"
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
	PingOk Ping = 1
)

type Node struct {
	ID 			int
	Name 		string 		`gorm:"name"`
	Addr 			string 		`gorm:"addr"`
	Type 		int8 		`gorm:"type"`
	Enable 		int8 		`gorm:"enable"`
	Ping 		Ping 		`gorm:"ping"`
	CreateAt  	time.Time	`gorm:"create_at"`
}

type NodeJson struct {
	ID 			int			`json:"id"`
	Name 		string 		`json:"name"`
	Addr 			string 		`json:"addr"`
	Type 		int8 		`json:"type"`
	Enable 		int8 		`json:"enable"`
	Ping 		Ping 		`json:"ping"`
	CreateAt  	string		`json:"create_at"`
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

	db := DPool().Conn
	db = db.Create(&node)
	if err = db.Error; err != nil {
		logger.Error("create err %v, exp %s .", err, db.QueryExpr())
		return errors.New("create node err")
	}

	return nil
}

func (node *Node) Del() bool {
	db := DPool().Conn
	db = db.Delete(node)
	if err := db.Error; err != nil {
		logger.Error("delete err %v, exp %s .", err, db.QueryExpr())
		return false
	}

	return true
}

func (node *Node) Update(data map[string]interface{}) error {
	db := DPool().Conn

	db = db.Model(&node).Update(data)
	if err := db.Error; err != nil {
		logger.Error("update err %v, exp %s .", err, db.QueryExpr())
		return errors.New("update node failed")
	}

	return nil
}

func (node *Node) FetchInfo() (Node, error) {
	var n Node
	db := DPool().Conn
	db = db.Where("id = ?", node.ID).First(&n)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		logger.Error("FetchInfo err %v, exp %s .", err, db.QueryExpr())
		return Node{}, errors.New("fetch node info failed")
	}

	return n, nil
}

func (node *Node) FetchRows(query string, args ...interface{}) ([]Node, error){
	db := DPool().Conn

	var list []Node
	db = db.Where(query, args...).Find(&list)
	if err := db.Error; err != nil {
		logger.Error("FetchRows err %v, exp %s .", err, db.QueryExpr())
		return nil, errors.New("fetchrows node failed")
	}

	return list, nil
}

func (node *Node) FormatJson() (NodeJson, error) {
	json := NodeJson{
		ID: node.ID,
		Name: node.Name,
		Addr: node.Addr,
		Type: node.Type,
		Enable: node.Enable,
		Ping: node.Ping,
		CreateAt: node.CreateAt.Format("2006-01-02 15:04:05"),
	}

	return json, nil
}