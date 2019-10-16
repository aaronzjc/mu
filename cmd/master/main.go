package main

import (
	"mu/internal/app/mu"
	"mu/internal/route"
	"mu/internal/util/logger"
)

func main() {
	//fmt.Println("fuck")
	//links, err := lib.FSite("tieba").BuildUrl()
	//fmt.Println(links)
	//if err != nil {
	//	fmt.Printf("%v", err)
	//}
	//for _, val := range links {
	//	fmt.Println(val.Url)
	//	page, err := val.Sp.CrawPage(val)
	//	if err != nil {
	//		fmt.Printf("%v", err)
	//	}
	//	fmt.Println(page.List)
	//}
	//
	//return

	// 注册路由
	route.RegisterRoutes()
	route.RegisterStatic()

	logger.Fatal(mu.App.Gin.Run(mu.App.Config.Server.Addr))
}
