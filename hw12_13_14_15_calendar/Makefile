BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
DOCKER_COMPOSE_PATH="./deployments/docker-compose.yaml"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

docker-compose-up:
	docker-compose -f $(DOCKER_COMPOSE_PATH) up -d
	while ! docker-compose -f ./deployments/docker-compose.yaml exec --user postgres db psql -c "select 'db ready!'"; do sleep 1; done;

build:
	docker-compose -f $(DOCKER_COMPOSE_PATH) build
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" .

run: build docker-compose-up migrate
	$(BIN) serve-http --config ./config/config.yaml

down:
	docker-compose -f $(DOCKER_COMPOSE_PATH) down

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	go version
	golangci-lint run ./...

install-migrator:
	(which goose > /dev/null) || go install github.com/pressly/goose/v3/cmd/goose@latest

migrate: install-migrator
	goose --dir ./migrations postgres "postgres://postgres:password@localhost:5435/postgres?sslmode=disable" up

install-protoc:
ifeq (, $(shell which protoc))
ifeq ($(shell uname -s),Darwin)
	brew install protobuf
else
	apt install -y protobuf-compiler
endif
endif
	(which protoc-gen-go > /dev/null) || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	(which protoc-gen-go-grpc > /dev/null) || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate: install-protoc
	protoc \
	--proto_path=./api/ --go_out=./internal/server/pb --go-grpc_out=./internal/server/pb \
	--grpc-gateway_out ./internal/server/pb --grpc-gateway_opt logtostderr=true,paths=import,generate_unbound_methods=true \
	--openapiv2_out ./gen/openapiv2 --openapiv2_opt logtostderr=true,generate_unbound_methods=true \
	api/*.proto
	go generate ./...

.PHONY: docker-compose-up build run down build-img run-img version test lint migrate generate
