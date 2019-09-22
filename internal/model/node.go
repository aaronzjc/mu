package model

import (
	"crawler/internal/app"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"
)

type Node struct {
	ID 			int
	Name 		string 		`gorm:"name"`
	Ip 			string 		`gorm:"ip"`
	Type 		int8 		`gorm:"type"`
	Enable 		int8 		`gorm:"enable"`
	Ping 		int8 		`gorm:"ping"`
	CreateAt  	time.Time	`gorm:"create_at"`
}

type NodeJson struct {
	ID 			int			`json:"id"`
	Name 		string 		`json:"name"`
	Ip 			string 		`json:"ip"`
	Type 		int8 		`json:"type"`
	Enable 		int8 		`json:"enable"`
	Ping 		int8 		`json:"ping"`
	CreateAt  	string		`json:"create_at"`
}

func (node *Node) TableName() string {
	return "node"
}

func (node *Node) CheckArgs() error {
	if node.Name == "" {
		return errors.New("参数为空")
	}
	if match, _ := regexp.MatchString("^(\\d+)\\.(\\d+)\\.(\\d+)\\.(\\d+)$", node.Ip); !match {
		return errors.New("IP不规范")
	}

	return nil
}

func (node *Node) Create() error {
	tmp, err := node.FetchInfo()
	if err != nil {
		return errors.New("create node err")
	}

	if tmp.ID > 0 {
		return errors.New(fmt.Sprintf("node with %s exists", node.Ip))
	}

	db := app.App.DB.Conn
	db = db.Create(&node)
	if err = db.Error; err != nil {
		log.Printf("[error] create err %v, exp %s \n", err, db.QueryExpr())
		return errors.New("create node err")
	}

	return nil
}

func (node *Node) Update(data map[string]interface{}) error {
	db := app.App.DB.Conn

	db = db.Model(&node).Update(data)
	if err := db.Error; err != nil {
		log.Printf("[error] update err %v, exp %s \n", err, db.QueryExpr())
		return errors.New("update node failed")
	}

	return nil
}

func (node *Node) FetchInfo() (Node, error) {
	var n Node
	db := app.App.DB.Conn
	db = db.Where("id = ?", node.ID).First(&n)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		log.Printf("[error] FetchInfo err %v, exp %s \n", err, db.QueryExpr())
		return Node{}, errors.New("fetch node info failed")
	}

	return n, nil
}

func (node *Node) FetchRows(query string, args ...interface{}) ([]Node, error){
	db := app.App.DB.Conn

	var list []Node
	db = db.Where(query, args...).Find(&list)
	if err := db.Error; err != nil {
		log.Printf("[error] FetchRows err %v, exp %s \n", err, db.QueryExpr())
		return nil, errors.New("fetchrows node failed")
	}

	return list, nil
}

func (node *Node) FormatJson() (NodeJson, error) {
	json := NodeJson{
		ID: node.ID,
		Name: node.Name,
		Ip: node.Ip,
		Type: node.Type,
		Enable: node.Enable,
		Ping: node.Ping,
		CreateAt: node.CreateAt.Format("2006-01-02 15:04:05"),
	}

	return json, nil
}