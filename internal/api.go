package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/infra/cache"
	"github.com/aaronzjc/mu/internal/infra/db"
	"github.com/aaronzjc/mu/internal/route"
	"github.com/aaronzjc/mu/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"gorm.io/gorm"
)

func SetupApi(ctx *cli.Context) error {
	var err error
	if ctx.String("config") == "" {
		return errors.New("invalid config option, use -h get full doc")
	}
	// 初始化项目配置
	configFile := ctx.String("config")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return errors.New("config file not found")
	}
	var conf *config.Config
	if conf, err = config.LoadConfig(configFile); err != nil {
		return err
	}
	// 初始化日志组件
	if err = logger.Setup(conf.Name, conf.Log.File); err != nil {
		return err
	}
	// 初始化DB
	if err = db.Setup(conf, &gorm.Config{}); err != nil {
		return err
	}
	// 初始化缓存
	if err = cache.Setup(&conf.Redis); err != nil {
		return err
	}

	// 初始化网站配置
	service.NewSiteService(store.NewSiteRepo(), nil).Init(context.Background())

	// 设置调试模式
	if conf.Env != "prod" {
		logger.SetLevel(conf.Log.Level)
		gin.SetMode(gin.DebugMode)
	}
	return nil
}

func RunApi(ctx *cli.Context) error {
	conf := config.Get()
	var addr string
	if conf.Env != "prod" {
		addr = fmt.Sprintf("127.0.0.1:%d", conf.Http.Port)
	} else {
		addr = fmt.Sprintf(":%d", conf.Http.Port)
	}
	// 启动服务器
	app := gin.New()
	route.Setup(app)
	server := &http.Server{
		Addr:         addr,
		Handler:      app,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}
	go server.ListenAndServe()
	logger.Info("[START] server listen at :", conf.Http.Port)

	// 监听关闭信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM)
	<-sig

	// 收到关闭信号，主动回收连接
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(ctxTimeout); err != nil {
		logger.Error("[STOP] server shutdown error", err)
		return err
	}
	logger.Info("[STOP] server shutdown ok")
	return nil
}
