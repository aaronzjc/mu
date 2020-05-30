GO111MODULE=on

.PHONY: production
production: vue clean agent commander api

agent:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/agent ./cmd/agent/main.go

commander:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/commander ./cmd/commander/main.go

api:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/api ./cmd/api/main.go

vue:
	cd web && npm install && npm run build

.PHONY: staging
staging: vue-staging

vue-staging:
	cd web && npm install && npm run build-staging

.PHONY: dev
dev: api-dev agent-dev commander-dev vue-dev vue-dev-build
api-dev:
	-rm ./bin/mu
	go fmt ./...
	go build -o ./bin/api ./cmd/api/main.go
	./bin/api

commander-dev:
	go build -o ./bin/commander ./cmd/commander/main.go
	./bin/commander

agent-dev:
	go build -o ./bin/agent ./cmd/agent/main.go
	./bin/agent

vue-dev:
	cd web && npm run serve

vue-dev-build:
	cd web && npm install && npm run build-dev

.PHONY: clean
clean:
	-rm ./bin/*