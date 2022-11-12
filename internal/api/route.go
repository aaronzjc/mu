package api

import (
	"os"

	"github.com/aaronzjc/mu/internal/api/handler"
	"github.com/aaronzjc/mu/internal/api/handler/admin"
	"github.com/aaronzjc/mu/internal/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoute(app *gin.Engine) {
	// 中间件
	app.Use(gin.Recovery(), gin.Logger())
	app.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// oauth
	rAuth := app.Group("/auth")
	{
		auth := handler.NewAuth()
		rAuth.GET("/config", auth.Platforms)
		rAuth.GET("/redirect", auth.Auth)
		rAuth.GET("/callback", auth.Callback)

		// 获取登录信息
		rAuth.Use(middleware.ApiAuth(false)).GET("/info/index", auth.LoginInfo)
		rAuth.Use(middleware.ApiAuth(true)).GET("/info/admin", auth.LoginInfo)
	}

	// index页面
	indexGroup := app.Group("/api")
	{
		idx := handler.NewIndex()
		indexGroup.GET("/sites", idx.Sites)
		indexGroup.GET("/news", idx.News)

		// 需要登陆的前端页面
		idxAuth := indexGroup.Use(middleware.ApiAuth(false))
		{
			// 收藏管理
			favor := handler.NewFavor()
			idxAuth.GET("/favors", favor.List)
			idxAuth.POST("/favors", favor.Add)
			idxAuth.POST("/favors/:id/del", favor.Remove)
		}
	}

	// admin页面
	adminGroup := app.Group("/admin").Use(middleware.ApiAuth(true))
	{
		// 节点管理
		node := admin.NewNode()
		adminGroup.GET("/nodes", node.List)
		adminGroup.POST("/nodes/:id/upsert", node.Upsert)
		adminGroup.GET("/nodes/:id/del", node.Del)

		// 站点管理
		site := admin.NewSite()
		adminGroup.GET("/sites", site.List)
		adminGroup.POST("/sites/:id/upsert", site.Upsert)
		adminGroup.POST("/sites/:id/craw", site.Craw)

		// 用户管理
		user := handler.NewUser()
		adminGroup.GET("/users", user.List)
	}

	// 静态资源托管
	RegistStatic(app)
}

func RegistStatic(r *gin.Engine) {
	r.Use(middleware.AddCacheControlHeader())
	path, _ := os.Getwd()
	dist := "/public/"
	r.StaticFile("/", path+dist+"index.html")
	r.StaticFile("/admin", path+dist+"admin.html")
	for _, v := range []string{"favicon.png", "index.manifest", "sw.js"} {
		r.StaticFile(v, path+dist+v)
	}
	for _, v := range []string{"pwa", "static"} {
		r.Static("/"+v, path+dist+""+v)
	}
}
