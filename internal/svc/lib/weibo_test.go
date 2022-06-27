package lib

import "testing"

func TestCrawWeibo(t *testing.T) {
	c := Weibo{NewSite(SITE_WEIBO)}
	links, _ := c.BuildUrl()
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:102.0) Gecko/20100101 Firefox/102.0"
	headers["Cookie"] = "SINAGLOBAL=2071023856171.5698.1622612955924; ULV=1656324033760:48:1:1:2691092481997.411.1656324033724:1644286815425; ALF=1671620901; UOR=,,mu.memosa.cn; SCF=ArUdQdJo1a0vp7QeaAG24sgclmBl55On9lp9C3EmxzqLbGVLvVBNl1QwD-9n_agTCl_CGr-KWiPlPJpN_t0ozM0.; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9WF--QHsfmEKLYfUnkcRW-V.; SUB=_2AkMVUF4_f8NxqwJRmP4RyWzlZIV_wgDEieKjDK_kJRMxHRl-yT9kqnActRB6PtBw0A7vr10ZHujj1KHS9SOzzZbsJiuV; _s_tentry=-; Apache=2691092481997.411.1656324033724"
	for _, link := range links {
		page, err := c.CrawPage(link, headers)
		if err != nil {
			t.Fatal("fetch weibo failed !")
		}
		t.Log(page.List)
	}
	t.Log("fetch weibo done .")
}