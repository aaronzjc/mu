package route

import (
	"crawler/internal/app"
	"crawler/internal/route/admin/node"
	"crawler/internal/route/admin/site"
	"crawler/internal/route/front"
	"github.com/gin-contrib/cors"
)

func RegisterStatic() {
	r := app.App.Gin

	path := "/Users/jincheng3/go/src/crawler"

	r.StaticFile("/", path + "/public/index.html")
	r.StaticFile("/admin", path + "/public/admin.html")

	r.StaticFile("favicon.png", path + "/public/favicon.png")
	r.Static("/static", path + "/public/")
}

func RegisterRoutes() {
	r := app.App.Gin

	c := cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
	})
	r.Use(c)

	// 前端路由
	r.GET("/aj", front.Aj)
	r.GET("/config", front.Config)

	// 后台管理路由
	api := r.Group("/api")
	{
		// 节点管理
		api.GET("/node", node.Info)
		api.GET("/node/list", node.List)
		api.POST("/node/upsert", node.CreateOrUpdateNode)

		// 站点管理
		api.GET("/site", site.Info)
		api.GET("/site/list", site.List)
		api.POST("/site/upsert", site.CreateOrUpdateNode)
	}
}
