package main

import (
	"crawler/internal/app"
	"crawler/internal/model"
	"crawler/internal/route"
	"log"
)

func main() {
	// 注册路由
	route.RegisterRoutes()
	route.RegisterStatic()

	// 初始化数据库
	(&model.Site{}).InitSites()

	log.Fatal(app.App.Gin.Run(app.App.Config.Addr))
}
