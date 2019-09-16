package main

import (
	"crawler/lib"
	"fmt"
	"github.com/robfig/cron"
	"log"
)

var linkPool = make(chan lib.Link, 3)
var pagePool = make(chan lib.Page, 3)

func addSite(s lib.Spider) {
	links, _ := s.BuildUrl()
	for _, link := range links {
		go func(link lib.Link) {
			linkPool <- link
		}(link)
	}
}

func addSites() {
	var spList []lib.Spider

	spList = append(spList, &lib.V2ex{
		Site: lib.NewSite(lib.SITE_V2EX),
	})
	spList = append(spList, &lib.Chouti{
		Site: lib.NewSite(lib.SITE_CT),
	})
	spList = append(spList, &lib.Zhihu{
		Site: lib.NewSite(lib.SITE_ZHIHU),
	})
	spList = append(spList, &lib.Weibo{
		Site: lib.NewSite(lib.SITE_WEIBO),
	})
	spList = append(spList, &lib.Hacker{
		Site: lib.NewSite(lib.SITE_HACKER),
	})

	for _, v := range spList {
		go addSite(v)
	}
}

func start() {
	for {
		select {
		case l := <-linkPool:
			go func() {
				sp := l.Sp
				page, err := sp.CrawPage(l)
				if err != nil {
					return
				}
				pagePool <- page
			}()
		case p := <-pagePool:
			go func() {
				sp := p.Link.Sp
				sp.Store(p)
			}()
		}
	}
}

func main() {
	addSites()
	start()
	return
	cron := cron.New()
	err := cron.AddFunc("0 */30 * * *", func() {
		fmt.Println("start crawling ...")
		addSites()
	})

	if err != nil {
		log.Fatal("[error] cron add err " + err.Error())
		return
	}
	cron.Start()

	start()
}
