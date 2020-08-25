package main

import (
	"log"
	"mu/internal/app/commander"
	"net/http"
	_ "net/http/pprof" // 必须，引入 pprof 模块
)

func main() {
	// 初始化
	commander.InitCommander()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	addr := ":7970"
	commander.RegisterRpcServer(addr)
}
