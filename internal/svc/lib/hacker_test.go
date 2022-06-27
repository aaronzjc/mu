package lib

import "testing"

func TestCrawHacker(t *testing.T) {
	c := Hacker{NewSite(SITE_HACKER)}
	links, _ := c.BuildUrl()
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:102.0) Gecko/20100101 Firefox/102.0"
	for _, link := range links {
		page, err := c.CrawPage(link, headers)
		if err != nil {
			t.Fatal("fetch hacker news failed !")
		}
		t.Log(page.List)
	}
	t.Log("fetch hacker news done .")
}