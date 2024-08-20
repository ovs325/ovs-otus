BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: all clean

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./

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
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.59.1



lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

install_prot: 
	sudo apt install protobuf-compiler
	go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/internal/genopenapi@latest


install_bin:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway 
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

compile:
	protoc -I . --go_out=./grpc --go_opt=paths=import --go-grpc_out=./grpc --go-grpc_opt=paths=import ./proto/event_service.proto
	protoc -I . --grpc-gateway_out=./grpc --grpc-gateway_opt=paths=import --grpc-gateway_opt=generate_unbound_methods=true ./proto/event_service.proto
	protoc -I . -I ./google/api/ --openapiv2_out ./ ./proto/event_service.proto

prot:
	protoc --proto_path=. --go-grpc_out=./grpc --go_out=./grpc --go-grpc_opt=paths=import --go_opt=paths=import ./proto/event_service.proto


postgr: 
	docker run --rm -d -e POSTGRES_PASSWORD=pass -e POSTGRES_USER=user -e POSTGRES_DB=calendar -p 7777:5432 postgres

rabbit: 
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management


all: calendar calendar_scheduler calendar_sender

calendar: cmd/calendar/main.go
	go build -o calendar cmd/calendar/main.go

calendar_scheduler: cmd/scheduler/main.go
	go build -o calendar_scheduler cmd/scheduler/main.go

calendar_sender: cmd/sender/main.go
	go build -o calendar_sender cmd/sender/main.go

clean:
	rm -f calendar calendar_scheduler calendar_sender

run-calendar:
	./calendar --config=./

run-scheduler:
	./calendar_scheduler --config=./

run-sender:
	./calendar_sender --config=./

.PHONY: run-calendar run-scheduler run-sender