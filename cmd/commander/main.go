package main

import (
	"mu/internal/app/commander"
)

func main() {
	// 初始化
	commander.InitCommander()

	addr := ":7970"
	commander.RegisterRpcServer(addr)
}
