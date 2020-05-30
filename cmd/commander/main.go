package main

import (
	"mu/internal/app/commander"
	"mu/internal/svc/schedule"
)

func initSchedule() {
	schedule.JobSchedule.InitJobs()
	schedule.JobSchedule.InitPool()
}

func init() {
	// 注册任务
	initSchedule()
}

func main() {
	addr := ":7970"
	commander.RegisterRpcServer(addr)
}
