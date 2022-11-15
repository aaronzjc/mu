package service

import (
	"context"
	"testing"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/application/store"
	"github.com/aaronzjc/mu/internal/core/site"
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCraw(t *testing.T) {
	require.Nil(t, test.SetupDb())
	require.Nil(t, test.SetupCache())

	nodeRepo := store.NewNodeRepo()
	siteRepo := store.NewSiteRepo()
	crawSvc := NewCrawService(siteRepo, nodeRepo)
	ctx := context.Background()

	ds := &dto.Site{
		Key:        site.SITE_GITHUB,
		NodeOption: model.ByHosts,
		NodeHosts:  []int{4},
	}
	// pick agent
	node, err := crawSvc.PickAgent(ctx, ds)
	require.Nil(t, err)
	require.Equal(t, node.ID, 4)

	// craw data
	err = crawSvc.Craw(ctx, ds)
	require.Nil(t, err)
	s := site.NewGithub()
	for _, tab := range s.Tabs {
		news, err := siteRepo.GetNews(ctx, s.Key, tab.Tag)
		require.Nil(t, err)
		assert.NotEmpty(t, news.List)
	}
}
