package mu

import (
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/util/config"
	"mu/internal/util/logger"
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

func (ins *Instance) initSites() {
	(&model.Site{}).InitSites()
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

	// 初始化站点
	App.initSites()
}
