package route

import (
	"crawler/internal/app/mu"
	"crawler/internal/route/admin/auth"
	"crawler/internal/route/admin/node"
	"crawler/internal/route/admin/site"
	"crawler/internal/route/index"
	"crawler/internal/route/middleware"
	"github.com/gin-contrib/cors"
	"os"
	"path/filepath"
)

func RegisterStatic() {
	r := mu.App.Gin

	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Dir(pwd)

	r.StaticFile("/", path + "/public/index.html")
	r.StaticFile("/admin", path + "/public/admin.html")
	r.StaticFile("/admin/login", path + "/public/login.html")

	r.StaticFile("favicon.png", path + "/public/favicon.png")
	r.Static("/static", path + "/public/static")
}

func RegisterRoutes() {
	r := mu.App.Gin

	c := cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	})
	r.Use(c)

	// 前端路由
	r.GET("/aj", index.Aj)
	r.GET("/config", index.Config)

	// 后台管理路由
	r.GET("/admin/auth", auth.Auth)
	r.GET("/admin/callback", auth.Callback)
	api := r.Group("/api")
	api.Use(middleware.Auth())
	{
		api.GET("/info", auth.Info)
		// 节点管理
		api.GET("/node", node.Info)
		api.GET("/node/list", node.List)
		api.POST("/node/upsert", node.CreateOrUpdateNode)

		// 站点管理
		api.GET("/site", site.Info)
		api.GET("/site/list", site.List)
		api.POST("/site/update", site.UpdateSite)
	}
}
