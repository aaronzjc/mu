package schedule

import (
	"context"
	"crawler/internal/model"
	"crawler/internal/svc/lib"
	"crawler/internal/svc/rpc"
	"crawler/internal/util/cache"
	"crawler/internal/util/logger"
	"encoding/json"
	"errors"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"math/rand"
	"sync"
	"time"
)

var (
	JobSchedule Schedule

	Pool = RpcPool{
		Clients: make(map[string]*RpcClient),
		Lock:    sync.RWMutex{},
	}
)

type RpcClient struct {
	Conn   *grpc.ClientConn
	Client *rpc.AgentClient
}

type RpcPool struct {
	Clients map[string]*RpcClient
	Lock    sync.RWMutex
}

func (r *RpcPool) Get(addr string) (*RpcClient, error) {
	r.Lock.RLock()
	rc, ok := r.Clients[addr]
	r.Lock.RUnlock()
	if ok {
		return rc, nil
	}

	client, err := r.Set(addr)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *RpcPool) Set(addr string) (*RpcClient, error) {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	rc, ok := r.Clients[addr]
	if ok {
		return rc, nil
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		//grpc.WithKeepaliveParams(keepalive.ClientParameters{
		//	Time:                20 * time.Second,
		//	Timeout:             3 * time.Second,
		//	PermitWithoutStream: true,
		//}),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logger.Error("connect error " + err.Error())
		return nil, errors.New("dial server " + addr + " failed")
	}

	client := rpc.NewAgentClient(conn)
	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	r.Clients[addr] = &RpcClient{
		Conn:   conn,
		Client: &client,
	}

	return r.Clients[addr], nil
}

func (r *RpcPool) Release(addr string) bool {
	r.Lock.Lock()
	rc, ok := r.Clients[addr]
	r.Lock.Unlock()
	if !ok {
		return true
	}

	delete(r.Clients, addr)
	_ = rc.Conn.Close()

	return true
}

/**
 * 抓取任务
 */
type CrawlerJob struct {
	Site model.Site
}

/**
 *	挑选一个节点执行。这里先简单粗暴一点，后面再思考怎么更加优雅的实现这块逻辑
 */
func (j *CrawlerJob) PickAgent() (model.Node, error) {
	var err error
	var idx int
	var nodes []model.Node
	rand.Seed(time.Now().UnixNano())
	if j.Site.NodeOption == model.ByType {
		nodes, err = (&model.Node{}).FetchRows("`type` = ? AND `ping` = ?", j.Site.NodeType, model.PingOk)
		if err != nil {
			logger.Error("pick agent error, err " + err.Error())
			return model.Node{}, errors.New("fetch nodes failed")
		}

	} else {
		var hosts []int
		err = json.Unmarshal([]byte(j.Site.NodeHosts), &hosts)
		if err != nil {
			return model.Node{}, errors.New("json unmarshal hosts err ")
		}
		if len(hosts) == 0 {
			return model.Node{}, errors.New("no available nodes")
		}
		nodes, err = (&model.Node{}).FetchRows("`id` IN (?) AND `enable` = ? AND `ping` = ?", hosts, model.Enable, model.PingOk)
		if err != nil {
			logger.Error("pick agent error, err " + err.Error())
			return model.Node{}, errors.New("fetch nodes failed")
		}
	}
	if len(nodes) == 0 {
		return model.Node{}, errors.New("no available nodes")
	}
	idx = rand.Int() % len(nodes)

	logger.Info("job [%s] pick agent [%s]", j.Site.Key, nodes[idx].Name)

	return nodes[idx], nil
}

/**
 * 必须是Run方法。实现Cron的Job接口。
 */
func (j *CrawlerJob) Run() {
	var err error
	node, err := j.PickAgent()

	if err != nil {
		logger.Error("pick agent error %v .", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	client, err := Pool.Get(node.Addr)
	if err != nil {
		logger.Error("get rpc client failed %v .", err)
		return
	}

	rpcClient := client.Client

	var result *rpc.Result
	result, err = (*rpcClient).Craw(ctx, &rpc.Job{Name: j.Site.Key})
	if err != nil {
		logger.Error("remote craw err %v", err)
		return
	}

	logger.Info("remote craw [%s] done", j.Site.Key)

	hotJson := new(lib.HotJson)
	hotJson.T = result.T
	for tag, hots := range result.Map {
		var list []lib.Hot
		for _, hot := range hots.Item {
			list = append(list, lib.Hot{
				Title:     hot.Title,
				Rank:      float64(hot.Rank),
				OriginUrl: hot.Url,
				Key: 	   hot.Key,
			})
		}
		hotJson.List = list
		data, err := json.Marshal(hotJson)
		if err != nil {
			logger.Error("Json_encode weibo error , err = %s .", err.Error())
			return
		}
		cache.SaveToRedis(j.Site.Key, tag, string(data))
	}

	return
}

func (j *CrawlerJob) ExecJobDirect() {
	site := lib.FSite(j.Site.Key)
	links, _ := site.BuildUrl()
	for _, link := range links {
		page, _ := link.Sp.CrawPage(link)
		hotJson := new(lib.HotJson)
		hotJson.T = time.Now().Format("2006-01-02 15:04:05")
		for _, hot := range page.List {
			hotJson.List = append(hotJson.List, hot)
		}
		hotJsonStr, _ := json.Marshal(hotJson)
		cache.SaveToRedis(j.Site.Key, link.Tag, string(hotJsonStr))
	}
}

/**
 * 服务存活检查任务
 */
type CheckJob struct {
	Name string
	Spec 	string
}

func (j *CheckJob) Run() {
	nodes, err := (&model.Node{}).FetchRows("`enable` = ?", model.Enable)
	if err != nil {
		panic("init pool failed " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, node := range nodes {
		client, err := Pool.Get(node.Addr)
		if err != nil {
			logger.Info("rpc health check : [%s] fetch conn error, err %v.", node.Name, err)
			if node.Ping != model.PingFailed {
				_ = node.Update(map[string]interface{}{
					"ping": model.PingFailed,
				})
			}
			continue
		}
		rpcClient := client.Client
		ping := &rpc.Ping{
			Ping: "ping",
		}
		result, err := (*rpcClient).Check(ctx, ping)
		if err != nil || result.Pong != ping.Ping {
			logger.Info("rpc health check : [%s] ping error, err %v.", node.Name, err)
			if node.Ping != model.PingFailed {
				_ = node.Update(map[string]interface{}{
					"ping": model.PingFailed,
				})
			}
			continue
		}

		logger.Info("rpc health check : [%s] is online.", node.Name)
		if node.Ping != model.PingOk {
			_ = node.Update(map[string]interface{}{
				"ping": model.PingOk,
			})
		}
	}
}

type Schedule struct {
	// 定时任务
	JobCron *cron.Cron

	// 全局变量
	JobMap sync.Map
}

func (s *Schedule) InitJobs() {
	m := model.Site{}
	sites, err := m.FetchRows("`enable` = ?", model.Enable)
	if err != nil {
		panic("schedule init sites failed " + err.Error())
	}

	for _, site := range sites {
		err = s.AddJob(site)
	}
}

func (s *Schedule) InitPool() {
	job := &CheckJob{
		Name: "heart_beat",
		Spec: "*/5 * * * *",
	}

	// 增加一个协程检查服务器状态
	_, _ = s.JobCron.AddJob(job.Spec, job)
}

func (s *Schedule) AddJob(site model.Site) error {
	job := &CrawlerJob{
		Site: site,
	}
	if _, ok := s.JobMap.Load(site.Key); ok {
		logger.Error("add job failed, job [%s] exists.", site.Key)
		return errors.New("job exists")
	}
	cronId, err := s.JobCron.AddJob(site.Cron, job)
	if err != nil {
		logger.Error("add job %s failed err = %v.", site.Key, err)
		return errors.New("add cron job failed")
	}

	// 将cron信息存储至全局的变量，方便管理维护
	s.JobMap.Store(site.Key, cronId)

	logger.Info("add job %s - [%s] success.", site.Key, site.Cron)

	return nil
}

func (s *Schedule) RemoveJob(siteKey string) bool {
	cronId, ok := s.JobMap.Load(siteKey)
	if !ok {
		logger.Info("job not exists in map")
		return true
	}
	s.JobCron.Remove(cronId.(cron.EntryID))
	s.JobMap.Delete(siteKey)

	logger.Info("remove job %s success .", siteKey)

	return true
}

func (s *Schedule) UpdateJob(site model.Site) bool {
	_, exist := s.JobMap.Load(site.Key)
	if exist {
		s.RemoveJob(site.Key)
	}

	err := s.AddJob(site)
	if err != nil {
		logger.Error("add job failed " + err.Error())
		return false
	}

	logger.Info("update job %s - [%s] success .", site.Key, site.Cron)

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
		if job, ok := entry.Job.(*CrawlerJob); ok {
			cm[job.Site.Key] = map[string]interface{}{
				"entry_id": entry.ID,
				"cron": job.Site.Cron,
				"next": next,
			}
			continue
		}
		if job, ok := entry.Job.(*CheckJob); ok {
			cm[job.Name] = map[string]interface{}{
				"entry_id": entry.ID,
				"cron": job.Spec,
				"next": next,
			}
			continue
		}
	}

	return map[string]interface{}{
		"JobMap": jm,
		"CronMap": cm,
	}
}

func init() {
	JobSchedule = Schedule{
		JobCron: cron.New(),
		JobMap:  sync.Map{},
	}
	JobSchedule.JobCron.Start()
}