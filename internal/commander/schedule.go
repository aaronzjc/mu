package commander

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/core/rpc"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/pkg/logger"

	"github.com/robfig/cron/v3"
)

var (
	JobSchedule Schedule

	Pool = rpc.RpcPool{
		Clients: make(map[string]*rpc.RpcClient),
		Lock:    sync.RWMutex{},
	}
)

type Schedule struct {
	// 定时任务
	JobCron *cron.Cron
	// 全局变量
	JobMap sync.Map
	svc    service.SiteService
}

func (s *Schedule) InitJobs() error {
	ctx := context.Background()
	sites, err := s.svc.Get(ctx, &dto.Query{
		Query: "`enable` = ?",
		Args:  []interface{}{model.Enable},
	})
	if err != nil {
		logger.Error("init jobs err " + err.Error())
		return err
	}

	for _, site := range sites {
		crawJob := NewCrawJob(site)
		if err := s.AddJob(site.Key, site.Cron, crawJob); err != nil {
			continue
		}
	}
	checkJob := NewCheckJob("heart_beat", "* * * * *")
	// 增加一个定时任务检查服务器状态
	err = s.AddJob(checkJob.Name, checkJob.Spec, checkJob)
	if err != nil {
		logger.Error("add check job err " + err.Error())
	}

	return nil
}

func (s *Schedule) TruncateJobs() error {
	s.JobMap.Range(func(k, v interface{}) bool {
		s.JobCron.Remove(v.(cron.EntryID))
		s.JobMap.Delete(k)
		return true
	})
	return nil
}

func (s *Schedule) AddJob(name string, spec string, job cron.Job) error {
	if _, ok := s.JobMap.Load(name); ok {
		logger.Error(fmt.Sprintf("add job failed, job [%s] exists.", name))
		return errors.New("job exists")
	}
	cronId, err := s.JobCron.AddJob(spec, job)
	if err != nil {
		logger.Error(fmt.Sprintf("add job %s failed err = %v.", name, err))
		return errors.New("add cron job failed")
	}
	// 将cron信息存储至全局的变量，方便管理维护
	s.JobMap.Store(name, cronId)
	logger.Info(fmt.Sprintf("add job %s - [%s] success.", name, spec))
	return nil
}

func (s *Schedule) RemoveJob(name string) bool {
	cronId, ok := s.JobMap.Load(name)
	if !ok {
		logger.Info("job not exists in map")
		return true
	}
	s.JobCron.Remove(cronId.(cron.EntryID))
	s.JobMap.Delete(name)

	logger.Info("remove job [" + name + "] success .")

	return true
}

func (s *Schedule) UpdateJob(name string, spec string, job cron.Job) bool {
	if _, exist := s.JobMap.Load(name); exist {
		s.RemoveJob(name)
	}
	if err := s.AddJob(name, spec, job); err != nil {
		logger.Error("add job failed " + err.Error())
		return false
	}
	logger.Info(fmt.Sprintf("update job %s - [%s] success .", name, spec))
	return true
}

func Debug() map[string]interface{} {
	jm := make(map[cron.EntryID]string)
	r := func(k interface{}, v interface{}) bool {
		jm[v.(cron.EntryID)] = k.(string)
		return true
	}
	JobSchedule.JobMap.Range(r)

	cm := make(map[string]interface{})
	for _, entry := range JobSchedule.JobCron.Entries() {
		next := entry.Schedule.Next(time.Now()).Format("2006-01-02 15:04:05")
		if job, ok := entry.Job.(*CrawJob); ok {
			cm[job.site.Key] = map[string]interface{}{
				"entry_id": entry.ID,
				"cron":     job.site.Cron,
				"next":     next,
			}
			continue
		}
		if job, ok := entry.Job.(*CheckJob); ok {
			cm[job.Name] = map[string]interface{}{
				"entry_id": entry.ID,
				"cron":     job.Spec,
				"next":     next,
			}
			continue
		}
	}

	return map[string]interface{}{
		"JobMap":  jm,
		"CronMap": cm,
	}
}

func SetupSchedule() {
	JobSchedule = Schedule{
		JobCron: cron.New(),
		JobMap:  sync.Map{},
		svc:     service.NewSiteService(store.NewSiteRepo(), nil),
	}
	JobSchedule.JobCron.Start()
}
