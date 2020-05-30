package main

import (
	"mu/internal/app/agent"
)

func main() {
	addr := ":7990"
	agent.RegisterRpcServer(addr)
}
