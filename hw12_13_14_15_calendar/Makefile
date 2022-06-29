BIN_FILENAME := calendar_otus_hw
BIN := "./bin/$(BIN_FILENAME)"
DOCKER_IMG="calendar:develop"
DOCKER_COMPOSE_PATH="./deployments/docker-compose.yaml"
DOCKER_COMPOSE_PATH_MS="./deployments/docker-compose.ms.yaml"
MODULE_PATH="github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X ${MODULE_PATH}/cmd.release="develop" -X ${MODULE_PATH}/cmd.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X ${MODULE_PATH}/cmd.gitHash=$(GIT_HASH)

up:
	docker-compose --env-file .env -f $(DOCKER_COMPOSE_PATH) -f $(DOCKER_COMPOSE_PATH_MS) up --build

down:
	docker-compose --env-file .env -f $(DOCKER_COMPOSE_PATH) -f $(DOCKER_COMPOSE_PATH_MS) down --remove-orphans

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

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

build-wo-docker-ms:
	docker-compose --env-file .env -f $(DOCKER_COMPOSE_PATH) build
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" .

up-wo-docker-ms:
	docker-compose --env-file .env -f $(DOCKER_COMPOSE_PATH) up -d
	while ! docker-compose --env-file .env -f $(DOCKER_COMPOSE_PATH) exec --user postgres db psql -c "select 'db ready!'" > /dev/null; do sleep 1; done;
	while ! curl -f -s http://localhost:15672 > /dev/null; do sleep 1; done;

install-migrator:
	(which goose > /dev/null) || go install github.com/pressly/goose/v3/cmd/goose@latest

migrate: install-migrator
	goose --dir ./migrations postgres "postgres://postgres:password@localhost:5435/postgres?sslmode=disable" up

run: build-wo-docker-ms up-wo-docker-ms migrate
	mkdir -p logs
	nohup $(BIN) serve-http --config ./config/calendar_config.yaml 0<&- &> ./logs/calendar.log &
	nohup $(BIN) sender --config ./config/sender_config.yaml 0<&- &> ./logs/sender.log &
	nohup $(BIN) scheduler --config ./config/scheduler_config.yaml 0<&- &> ./logs/scheduler.log &

stop:
	killall $(BIN_FILENAME)
	docker-compose --env-file .env -f $(DOCKER_COMPOSE_PATH) down

version: build-wo-docker-ms
	$(BIN) version

.PHONY: up down test lint generate build-wo-docker-ms up-wo-docker-ms migrate run stop version
