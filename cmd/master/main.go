package main

import (
	"log"
	"os"
	"strings"

	"crawler/internal/route/front"
	"crawler/internal/util/config"

	"github.com/gin-gonic/gin"
)

func init() {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "production" || env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	appConfig := config.NewConfig()

	r := gin.New()

	path := "/Users/jincheng3/go/src/crawler"

	r.StaticFile("/", path + "/fakedist/index.html")
	r.StaticFile("/admin", path + "/fakedist/admin.html")

	r.StaticFile("favicon.png", path + "/fakedist/favicon.png")
	r.Static("/static", path + "/fakedist/static")

	r.GET("/aj", front.Aj)
	r.GET("/config", front.Config)

	log.Fatal(r.Run(appConfig.Addr))
}
