package main

import (
	"crawler/internal/app/mu"
	"crawler/internal/route"
	"crawler/internal/util/logger"
)

func main() {
	// 注册路由
	route.RegisterRoutes()
	route.RegisterStatic()

	logger.Fatal(mu.App.Gin.Run(mu.App.Config.Server.Addr))
}
