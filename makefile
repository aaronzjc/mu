GO111MODULE=on

.PHONY: crawler
crawler:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/crawler ./cron/main.go

.PHONY: mu
mu:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/mu ./server/main.go


