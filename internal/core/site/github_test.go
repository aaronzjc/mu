package site

import (
	"testing"

	"github.com/aaronzjc/mu/test"
	"github.com/stretchr/testify/assert"
)

func TestCrawGithub(t *testing.T) {
	test.SetupProxy()
	defer test.ClearProxy()

	assert := assert.New(t)

	c := &Github{
		Site{
			Name:     "Github",
			Key:      SITE_GITHUB,
			Root:     "https://github.com",
			Desc:     "Github.com",
			CrawType: CrawHtml,
			Tabs:     GithubTabs,
		},
	}
	links, _ := c.BuildUrl()
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:102.0) Gecko/20100101 Firefox/102.0"
	for _, link := range links {
		page, err := c.CrawPage(link, headers)
		assert.Nil(err)
		assert.NotEmpty(page.List)
	}
	t.Log("fetch github done .")
}
