package main

import (
	"mu/internal/app/crawler"
)

func main() {
	addr := ":7990"
	crawler.RegisterRpcServer(addr)
}
