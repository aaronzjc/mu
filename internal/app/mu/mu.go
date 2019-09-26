package mu

import (
	"crawler/internal/model"
	"crawler/internal/svc/schedule"
	"crawler/internal/util/config"
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
}

func (ins *Instance) initConfig() {
	defer log.Printf("[info] init config complete.\n")
	ins.Config = config.NewConfig()
}

func (ins *Instance) initCron() {
	schedule.JobSchedule.InitJobs()
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

	// 初始化数据库
	(&model.Site{}).InitSites()
}