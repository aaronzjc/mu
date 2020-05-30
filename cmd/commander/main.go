package main

import (
	"mu/internal/app/commander"
	"mu/internal/svc/schedule"
)

func init() {
	// 注册定时任务
	schedule.JobSchedule.InitJobs()
}

func main() {
	addr := ":7970"
	commander.RegisterRpcServer(addr)
}
