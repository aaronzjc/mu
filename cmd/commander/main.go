package main

import (
	"github.com/gin-gonic/gin"
	"mu/internal/app/commander"
	"mu/internal/svc/schedule"
	"strings"
	"os"
)

func init() {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "production" || env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// 注册定时任务
	schedule.JobSchedule.InitJobs()
}

func main() {
	addr := ":7970"
	commander.RegisterRpcServer(addr)
}
