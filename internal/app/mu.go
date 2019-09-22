package app

import (
	"crawler/internal/util/config"
	"crawler/internal/util/db"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

var (
	// App 程序实例
	App *Instance
)

type Instance struct {
	Gin		*gin.Engine
	Config  config.Config
	DB 		*db.DB
}

func (ins *Instance) initConfig() {
	defer log.Printf("[info] init config complete.\n")
	ins.Config = config.NewConfig()
}

func (ins *Instance) initDb() {
	defer log.Printf("[info] init db complete.\n")
	ins.DB = &db.DB{}
	ins.DB.Connect(&App.Config)
}

func init() {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "production" || env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	App = &Instance{
		Gin: gin.New(),
	}

	App.initConfig()
	App.initDb()
}