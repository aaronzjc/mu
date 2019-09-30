package mu

import (
	"crawler/internal/model"
	"crawler/internal/svc/schedule"
	"crawler/internal/util/config"
	"crawler/internal/util/logger"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

var (
	// App 程序实例
	App *Instance
)

type Instance struct {
	Gin    *gin.Engine
	Config config.Config
}

func (ins *Instance) initConfig() {
	defer logger.Info("init config complete.")
	ins.Config = config.NewConfig()
}

func (ins *Instance) initSchedule() {
	schedule.JobSchedule.InitJobs()
	schedule.JobSchedule.InitPool()
}

func init() {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "production" || env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	App = &Instance{
		Gin: gin.New(),
	}

	// 初始化配置
	App.initConfig()

	// 初始化数据库
	(&model.Site{}).InitSites()

	// 初始化任务队列，rpc等
	App.initSchedule()
}
