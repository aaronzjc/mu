GO111MODULE=on
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64

.PHONY: crawler
crawler:
	go build -o bin/crawler ./cron/main.go

.PHONY: mu
mu:
	go build -o bin/mu ./server/main.go

.PHONY: clean
clean:
	-rm ./bin/*