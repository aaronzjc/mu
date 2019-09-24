package main

import (
	"crawler/internal/app"
	"crawler/internal/route"
	"log"
)

func main() {
	// 注册路由
	route.RegisterRoutes()
	route.RegisterStatic()

	log.Fatal(app.App.Gin.Run(app.App.Config.Addr))
}
