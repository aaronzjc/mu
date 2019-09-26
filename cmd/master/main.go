package main

import (
	"crawler/internal/app/mu"
	"crawler/internal/route"
	"log"
)

func main() {
	// 注册路由
	route.RegisterRoutes()
	route.RegisterStatic()

	log.Fatal(mu.App.Gin.Run(mu.App.Config.Server.Addr))
}
