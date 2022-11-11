package dto

import (
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/pkg/helper"
)

type Node struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Addr     string `json:"addr"`
	Type     int8   `json:"type"`
	Enable   int8   `json:"enable"`
	Ping     int8   `json:"ping"`
	CreateAt string `json:"create_at"`
}

func (n *Node) FillByModel(node model.Node) *Node {
	n.ID = node.ID
	n.Name = node.Name
	n.Addr = node.Addr
	n.Type = node.Type
	n.Enable = node.Enable
	n.Ping = node.Ping
	n.CreateAt = helper.TimeToLocalStr(node.CreateAt)
	return n
}
