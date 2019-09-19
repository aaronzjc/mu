GO111MODULE=on

.PHONY: crawler
crawler:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/crawler ./cmd/node/main.go

.PHONY: mu
mu:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/mu ./cmd/master/main.go

.PHONY: vue
vue:
	cd web && npm run build

.PHONY: clean
clean:
	-rm ./bin/*