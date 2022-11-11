package commander

import (
	"context"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/service"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/pkg/logger"
	"github.com/robfig/cron/v3"
)

type CrawJob struct {
	site *dto.Site

	svc service.CrawService
}

var _ cron.Job = &CrawJob{}

func (j *CrawJob) Run() {
	ctx := context.Background()
	err := j.svc.Craw(ctx, j.site)
	if err != nil {
		logger.Error("craw job run err " + err.Error())
	}
}

func NewCrawJob(site *dto.Site) *CrawJob {
	return &CrawJob{
		site: site,
		svc:  service.NewCrawService(),
	}
}

/**
 * 服务存活检查任务
 */
type CheckJob struct {
	Name string
	Spec string

	svc service.NodeService
}

var _ cron.Job = &CheckJob{}

func (j *CheckJob) Run() {
	ctx := context.Background()
	j.svc.CheckNodes(ctx, &Pool)
}

func NewCheckJob(name string, spec string) *CheckJob {
	repo := store.NewNodeRepo()
	return &CheckJob{
		Name: name,
		Spec: spec,
		svc:  service.NewNodeService(repo),
	}
}
