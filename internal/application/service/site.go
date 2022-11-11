package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/constant"
	"github.com/aaronzjc/mu/internal/core/rpc"
	"github.com/aaronzjc/mu/internal/core/site"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/internal/domain/repo"
	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/pkg/logger"
)

type SiteService interface {
	Init(context.Context) error
	ListOfIndex(context.Context) ([]*dto.IndexSite, error)
	News(context.Context, string, string) (*dto.News, error)

	Upsert(context.Context, *dto.Site) error
	Del(context.Context, int) error
	Get(context.Context, *dto.Query) ([]*dto.Site, error)
}

type SiteServiceImpl struct {
	repo      repo.SiteRepo
	favorRepo repo.FavorRepo
}

var _ SiteService = &SiteServiceImpl{}

func (s *SiteServiceImpl) Init(ctx context.Context) error {
	var tagStr []byte

	for siteKey, spider := range site.SiteMap {
		siteConf := spider.GetSite()
		sites, err := s.repo.Get(ctx, &dto.Query{
			Query: " `key` = ? ",
			Args:  []interface{}{siteKey},
		})
		if err != nil {
			logger.Error("init site " + siteKey + " failed")
			continue
		}

		tagMap := make(map[string]int)
		if len(sites) > 0 {
			var siteTags []dto.Tag
			json.Unmarshal([]byte(sites[0].Tags), &siteTags)
			for _, v := range siteTags {
				if v.Enable != 0 {
					continue
				}
				tagMap[v.Key] = 1
			}
		}
		var tags []dto.Tag
		var e int8
		for _, tag := range siteConf.Tabs {
			e = 1
			if _, ok := tagMap[tag["tag"]]; ok {
				e = 0
			}
			tags = append(tags, dto.Tag{
				Key:    tag["tag"],
				Name:   tag["name"],
				Enable: e,
			})
		}
		tagStr, _ = json.Marshal(tags)

		if len(sites) > 0 {
			err = s.repo.Update(ctx, sites[0], map[string]interface{}{
				"name": siteConf.Name,
				"root": siteConf.Root,
				"tags": string(tagStr),
				"type": siteConf.CrawType,
			})
			if err != nil {
				logger.Error("update site err")
			}
			continue
		}
		err = s.repo.Create(ctx, model.Site{
			Name:       siteConf.Name,
			Key:        siteConf.Key,
			Root:       siteConf.Root,
			Cron:       "*/30 * * * *",
			NodeOption: 1, // 默认使用分类
			NodeType:   1, // 默认国内的机器
			NodeHosts:  "",
			Desc:       siteConf.Desc,
			Tags:       string(tagStr),
			Type:       siteConf.CrawType,
			ReqHeaders: "",
		})
		if err != nil {
			continue
		}
	}
	return nil
}

func (s *SiteServiceImpl) ListOfIndex(ctx context.Context) ([]*dto.IndexSite, error) {
	var idxSites []*dto.IndexSite
	sites, _ := s.repo.Get(ctx, &dto.Query{
		Query: "`enable` = ?",
		Args:  []interface{}{model.Enable},
	})
	for _, site := range sites {
		var tags []dto.Tag
		if site.Tags == "" {
			continue
		}
		idxSite := &dto.IndexSite{
			Name: site.Name,
			Key:  site.Key,
		}
		if err := json.Unmarshal([]byte(site.Tags), &tags); err != nil {
			continue
		}
		for _, tag := range tags {
			if tag.Enable == 0 {
				continue
			}
			idxSite.Tags = append(idxSite.Tags, dto.Tag{
				Name: tag.Name,
				Key:  tag.Key,
			})
		}
		idxSites = append(idxSites, idxSite)
	}
	return idxSites, nil
}

func (s *SiteServiceImpl) News(ctx context.Context, k string, kk string) (*dto.News, error) {
	var sn dto.News
	news, err := s.repo.GetNews(ctx, k, kk)
	if err != nil {
		return nil, err
	}
	favors, _ := s.favorRepo.Get(ctx, &dto.Query{
		Query: "`user_id` = ? AND `site` = ?",
		Args:  []interface{}{0, k},
	})
	mp := make(map[string]bool)
	for _, v := range favors {
		mp[v.Key] = true
	}
	sn.T = news.T
	for _, v := range news.List {
		_, ok := mp[v.Key]
		sn.List = append(sn.List, dto.NewsItem{
			Key:       v.Key,
			Title:     v.Title,
			Desc:      v.Desc,
			Rank:      v.Rank,
			OriginUrl: v.OriginUrl,
			Card:      v.Card,
			Ext:       v.Ext,
			Mark:      ok,
		})
	}
	return &sn, nil
}

func (s *SiteServiceImpl) Get(ctx context.Context, q *dto.Query) ([]*dto.Site, error) {
	var sites []*dto.Site
	mns, err := s.repo.Get(ctx, q)
	if err != nil {
		return sites, err
	}
	for _, v := range mns {
		sites = append(sites, (&dto.Site{}).FillByModel(v))
	}
	return sites, nil
}

func (s *SiteServiceImpl) Upsert(ctx context.Context, site *dto.Site) error {
	sites, err := s.Get(ctx, &dto.Query{
		Query: "`id` = ?",
		Args:  []interface{}{site.ID},
	})
	if err != nil {
		return err
	}
	tagBytes, _ := json.Marshal(site.Tags)
	hostsBytes, _ := json.Marshal(site.NodeHosts)
	reqHeaders, _ := json.Marshal(site.ReqHeaders)
	if len(sites) > 0 {
		s.repo.Update(ctx, model.Site{ID: sites[0].ID}, map[string]interface{}{
			"name":        site.Name,
			"key":         site.Key,
			"desc":        site.Desc,
			"cron":        site.Cron,
			"enable":      site.Enable,
			"node_option": site.NodeOption,
			"node_type":   site.NodeType,
			"tags":        string(tagBytes),
			"node_hosts":  string(hostsBytes),
			"req_headers": string(reqHeaders),
		})
	} else {
		s.repo.Create(ctx, model.Site{
			ID:         site.ID,
			Name:       site.Name,
			Key:        site.Key,
			Desc:       site.Desc,
			Cron:       site.Cron,
			Enable:     site.Enable,
			NodeOption: site.NodeOption,
			NodeType:   site.NodeType,
			Tags:       string(tagBytes),
			NodeHosts:  string(hostsBytes),
			ReqHeaders: string(reqHeaders),
		})
	}

	// 更新任务配置
	timeoutCtx, calcel := context.WithTimeout(ctx, time.Second*3)
	defer calcel()
	svcUrl := config.Get().GetServiceUrl(constant.SvcCommander)
	if svcUrl != "" {
		client, _ := rpc.NewCommanderClient(svcUrl)
		client.UpdateCron(timeoutCtx, &pb.Cron{Site: site.Key})
	} else {
		logger.Error("svc commander not found")
	}

	return nil
}

func (s *SiteServiceImpl) Del(ctx context.Context, id int) error {
	return s.repo.Del(ctx, model.Site{ID: id})
}

func NewSiteService(repo repo.SiteRepo, favor repo.FavorRepo) *SiteServiceImpl {
	return &SiteServiceImpl{
		repo:      repo,
		favorRepo: favor,
	}
}
