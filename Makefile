GO111MODULE=on

.PHONY: all
all: vue clean crawler mu

.PHONY: crawler mu vue
crawler:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/crawler ./cmd/node/main.go

mu:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/mu ./cmd/master/main.go

vue:
	cd web && npm install && npm run build

.PHONY: mu-dev mu-staging vue-staging
mu-dev:
	go build -o ./bin/mu ./cmd/master/main.go
	./bin/mu

mu-staging:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/mu ./cmd/master/main.go

vue-staging:
	cd web && mv .env.staging .env && npm install && npm run build

.PHONY: clean
clean:
	-rm ./bin/*