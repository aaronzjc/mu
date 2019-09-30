package main

import (
	"crawler/internal/app/crawler"
)

func main() {
	addr := ":7990"
	crawler.RegisterRpcServer(addr)
}
