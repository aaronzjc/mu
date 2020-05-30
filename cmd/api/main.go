package main

import (
	"mu/internal/app/mu"
	"mu/internal/route"
	"mu/internal/util/logger"
)

func main() {
	// 注册路由
	route.RegisterRoutes()

	logger.Fatal(mu.App.Gin.Run(mu.App.Config.Server.Port))
}
