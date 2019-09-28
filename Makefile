GO111MODULE=on

.PHONY: production
production: vue clean crawler mu

crawler:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/crawler ./cmd/node/main.go

mu:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/mu ./cmd/master/main.go

vue:
	cd web && npm install && npm run build

.PHONY: staging
staging: vue-staging

vue-staging:
	cd web && npm install && npm run build-staging

.PHONY: dev
dev: mu-dev crawler-dev vue-dev
mu-dev:
	go build -o ./bin/mu ./cmd/master/main.go
	./bin/mu

crawler-dev:
	go build -o ./bin/crawler ./cmd/node/main.go
	./bin/crawler

vue-dev:
	cd web && npm run serve

.PHONY: clean
clean:
	-rm ./bin/*