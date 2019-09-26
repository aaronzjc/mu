package main

import (
	"crawler/internal/app/crawler"
	"fmt"
)

func main() {
	addr := ":7990"
	crawler.RegisterRpcServer(addr)
	fmt.Println("server is listening on :7990")
}
