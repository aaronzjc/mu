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
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	r.Use(c)

	// oauth相关
	rAuth := r.Group("/oauth")
	{
		rAuth.GET("/config", oauth.Config)
		rAuth.GET("/auth", oauth.Auth)
		rAuth.GET("/callback", oauth.Callback)
	}

	// 前端路由
	api := r.Group("/api")
	{
		api.GET("/config", hot.Tabs)
		api.GET("/list", hot.List)

		apiAuth := api.Group("").Use(middleware.ApiAuth(false))
		{
			// 获取登录用户信息
			apiAuth.GET("/info", idxAuth.Info)
			// 收藏管理
			apiAuth.GET("/favor/list", favor.List)
			apiAuth.POST("/favor/add", favor.Add)
			apiAuth.POST("/favor/remove", favor.Remove)
		}
	}

	// 后台路由管理
	admin := r.Group("/admin")
	admin.Use(middleware.ApiAuth(true))
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

	// 如果后端托管，则注册静态资源路由
	cnf := config.NewConfig()
	if cnf.Server.Static {
		RegisterStatic()
	}
}
