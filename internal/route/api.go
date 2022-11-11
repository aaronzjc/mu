package route

import (
	"os"

	"github.com/aaronzjc/mu/internal/api"
	adminApi "github.com/aaronzjc/mu/internal/api/admin"
	"github.com/aaronzjc/mu/internal/route/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
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
		auth := api.NewAuth()
		rAuth.GET("/config", auth.Platforms)
		rAuth.GET("/redirect", auth.Auth)
		rAuth.GET("/callback", auth.Callback)

		// 获取登录信息
		rAuth.Use(middleware.ApiAuth(false)).GET("/info/index", auth.LoginInfo)
		rAuth.Use(middleware.ApiAuth(true)).GET("/info/admin", auth.LoginInfo)
	}

	// index页面
	index := app.Group("/api")
	{
		idx := api.NewIndex()
		index.GET("/sites", idx.Sites)
		index.GET("/news", idx.News)

		// 需要登陆的前端页面
		idxAuth := index.Use(middleware.ApiAuth(false))
		{
			// 收藏管理
			favor := api.NewFavor()
			idxAuth.GET("/favors", favor.List)
			idxAuth.POST("/favors", favor.Add)
			idxAuth.POST("/favors/:id/del", favor.Remove)
		}
	}

	// admin页面
	admin := app.Group("/admin").Use(middleware.ApiAuth(true))
	{
		// 节点管理
		node := adminApi.NewNode()
		admin.GET("/nodes", node.List)
		admin.POST("/nodes/:id/upsert", node.Upsert)
		admin.GET("/nodes/:id/del", node.Del)

		// 站点管理
		site := adminApi.NewSite()
		admin.GET("/sites", site.List)
		admin.POST("/sites/:id/upsert", site.Upsert)
		admin.POST("/sites/:id/craw", site.Craw)

		// 用户管理
		user := api.NewUser()
		admin.GET("/users", user.List)
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
