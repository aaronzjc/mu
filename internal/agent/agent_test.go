package agent

import (
	"context"
	"testing"

	"github.com/aaronzjc/mu/internal/core/site"
	"github.com/aaronzjc/mu/internal/pb"
	"github.com/aaronzjc/mu/test"
	"github.com/stretchr/testify/assert"
)

var (
	agent = NewAgentServer()
)

func TestCraw(t *testing.T) {
	test.SetupProxy()
	defer test.ClearProxy()

	cases := []site.Site{
		site.NewGithub().Site,
		site.NewHacker().Site,
	}
	for _, v := range cases {
		result, err := agent.Craw(context.Background(), &pb.Job{
			Name: v.Key,
		})
		assert.Nil(t, err)
		for _, vv := range v.Tabs {
			tabStr, ok := result.HotMap[vv.Tag]
			assert.True(t, ok)
			assert.True(t, tabStr != "[]")
		}
	}
}
