GO111MODULE=on
GOOS=linux

.PHONY: crawler
crawler:
	go build -o bin/crawler ./cron/main.go

.PHONY: mu
mu:
	go build -o bin/mu ./server/main.go


