package route

import (
	"github.com/gin-contrib/cors"
	"mu/internal/app/mu"
	adminAuth "mu/internal/route/admin/auth"
	"mu/internal/route/admin/node"
	"mu/internal/route/admin/site"
	"mu/internal/route/admin/user"
	idxAuth "mu/internal/route/index/auth"
	"mu/internal/route/index/favor"
	"mu/internal/route/index/hot"
	"mu/internal/route/middleware"
	"mu/internal/route/oauth"
	"mu/internal/util/config"
	"os"
	"path/filepath"
)

func RegisterStatic() {
	r := mu.App.Gin

	r.Use(middleware.AddCacheControlHeader())

	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Dir(pwd)

	r.StaticFile("/", path+"/public/index.html")
	r.StaticFile("/admin", path+"/public/admin.html")

	for _, v := range []string{"favicon.png", "index.manifest", "sw.js"} {
		r.StaticFile(v, path+"/public/"+v)
	}

	for _, v := range []string{"pwa", "static"} {
		r.Static("/"+v, path+"/public/"+v)
	}
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

	// Auth操作
	r.GET("/auth_config", oauth.Config)
	r.GET("/oauth/auth", oauth.Auth)
	r.GET("/oauth/callback", oauth.Callback)

	// 前端路由
	r.GET("/config", hot.Tabs)
	r.GET("/logout", idxAuth.Logout)
	idx := r.Group("")
	idx.Use(middleware.ApiAuth(false))
	{
		// 本组路由获取用户信息，但是不强制登录
		idx.GET("/info", idxAuth.Info)
		idx.GET("/list", hot.List)
	}
	api := r.Group("/api")
	api.Use(middleware.ApiAuth(true))
	{
		// 收藏管理
		api.GET("/favor/list", favor.List)
		api.POST("/favor/add", favor.Add)
		api.POST("/favor/remove", favor.Remove)
	}

	// 后台路由管理
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuth())
	{
		admin.GET("/debug", site.Debug)
		admin.GET("/info", adminAuth.Info)
		// 节点管理
		admin.GET("/node", node.Info)
		admin.GET("/node/list", node.List)
		admin.POST("/node/upsert", node.CreateOrUpdateNode)
		admin.GET("/node/del", node.Del)

		// 站点管理
		admin.GET("/site", site.Info)
		admin.GET("/site/list", site.List)
		admin.POST("/site/update", site.UpdateSite)

		// 用户管理
		admin.GET("/user/list", user.List)
	}

	// 如果配置里面没有配置前端URL，这里后端托管
	cnf := config.NewConfig()
	if cnf.WebUrl() == "" {
		RegisterStatic()
	}
}
