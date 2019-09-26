package schedule

import (
	"context"
	"crawler/internal/model"
	"crawler/internal/svc/lib"
	"crawler/internal/svc/rpc"
	"crawler/internal/util/cache"
	"encoding/json"
	"errors"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	JobSchedule Schedule
)

type Schedule struct {
	// 定时任务
	JobCron	*cron.Cron

	// 全局变量
	JobMap 	sync.Map
}

func init() {
	JobSchedule = Schedule{
		JobCron: cron.New(),
		JobMap: sync.Map{},
	}
}

func (s *Schedule) InitJobs() {
	m := model.Site{}
	sites, err := m.FetchRows("`enable` = ?", 1)
	if err != nil {
		panic("schedule init sites failed " + err.Error())
	}

	for _, site := range sites {
		err = s.AddJob(&site)
	}
}

func (s *Schedule) AddJob(site *model.Site) error {
	job := &Job{
		Site: site,
	}
	cronId, err := s.JobCron.AddFunc(site.Cron, job.ExecJob)
	if err != nil {
		log.Printf("[error] add job %s failed err = %v.\n", site.Name, err)
		return errors.New("add cron job failed")
	}

	// 将cron信息存储至全局的变量，方便管理维护
	s.JobMap.Store(site.Name, cronId)

	log.Printf("[info] add job %s - [%s] success\n", site.Name, site.Cron)

	return nil
}

func (s *Schedule) RemoveJob(siteName string) bool {
	cronId, ok := s.JobMap.Load(siteName)
	if !ok {
		log.Printf("[warning] job not exists in map")
		return true
	}
	s.JobCron.Remove(cronId.(cron.EntryID))
	s.JobMap.Delete(cronId)

	log.Printf("[info] remove job %s success .\n", siteName)

	return true
}

func (s *Schedule) UpdateJob(site *model.Site) bool {
	_, exist := s.JobMap.Load(site.Name)
	if exist {
		s.RemoveJob(site.Name)
	}

	err := s.AddJob(site)
	if err != nil {
		log.Printf("[error] add job failed " + err.Error())
		return false
	}

	log.Printf("[info] update job %s - [%s] success .\n", site.Name, site.Cron)

	return true
}

type Job struct {
	Site *model.Site
}

/**
 *	挑选一个节点执行。这里先简单粗暴一点，后面再思考怎么更加优雅的实现这块逻辑
 */
func (j *Job) PickAgent() (model.Node, error) {
	var err error
	var idx int
	rand.Seed(time.Now().UnixNano())
	if j.Site.NodeOption == model.ByType {
		nodes, err := (&model.Node{}).FetchRows("`type` = ? AND `ping` = ?", j.Site.NodeType, model.PingOk)
		if err != nil {
			log.Printf("[error] pick agent error, err " + err.Error())
			return model.Node{}, errors.New("fetch nodes failed")
		}
		if len(nodes) == 0 {
			return model.Node{}, errors.New("no available nodes")
		}
		idx = rand.Int() % len(nodes)
		return nodes[idx], nil
	} else {
		var hosts []int
		err = json.Unmarshal([]byte(j.Site.NodeHosts), &hosts)
		if err != nil {
			return model.Node{}, errors.New("json unmarshal hosts err ")
		}
		if len(hosts) == 0 {
			return model.Node{}, errors.New("no available nodes")
		}

		idx = rand.Int() % len(hosts)
		node, err := (&model.Node{
			ID: hosts[idx],
		}).FetchInfo()
		if err != nil {
			log.Printf("[error] pick agent error, err " + err.Error())
			return model.Node{}, errors.New("fetch nodes failed")
		}

		return node, nil
	}
}

func (j *Job) ExecJob() {
	var err error
	node, err := j.PickAgent()

	if err != nil {
		log.Printf("[error] pick agent error %v \n", err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(node.Addr, opts...)
	if err != nil {
		log.Fatal("[error] connect error " + err.Error())
	}
	defer conn.Close()

	client := rpc.NewAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()

	var result *rpc.Result
	result, err = client.Craw(ctx, &rpc.Job{Name: j.Site.Key})
	if err != nil {
		log.Printf("[error] remote craw err %v", err)
		return
	}

	hotJson := new(lib.HotJson)
	hotJson.T = result.T
	for tag, hots := range result.Map {
		var list []lib.Hot
		for _, hot := range hots.Item {
			list = append(list, lib.Hot{
				Title: hot.Title,
				Rank: float64(hot.Rank),
				OriginUrl: hot.Url,
			})
		}
		hotJson.List = list
		data, err := json.Marshal(hotJson)
		if err != nil {
			log.Printf("[error] Json_encode weibo error , err = %s\n", err.Error())
			return
		}
		cache.SaveToRedis(j.Site.Name, tag, string(data))
	}

	return
}