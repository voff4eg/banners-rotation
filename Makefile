ifeq ($(OS),Windows_NT)
    GOOS := windows
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        GOOS := linux
    endif
    ifeq ($(UNAME_S),Darwin)
        GOOS := darwin
    endif
endif

.PHONY: dev-build
dev-build:
	go build -v ./cmd/server.go

.PHONY: build
build:
	GOOS=${GOOS} CGO_ENABLED=0 GOARCH=amd64 go build -v -trimpath ./cmd/server.go

.PHONY: test-unit
test-unit:
	go test ./internal/... -v -race -count=10

.PHONY: test-integration
test-integration:
	go test ./test/integration/... -v

.PHONY: test
test: test-unit test-integration

.PHONY: run
run: docker-build docker-up

.PHONY: lint
lint:
	golangci-lint run

.PHONY: docker-build
docker-build:
	docker-compose -f docker-compose.yaml build

.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yaml up -d --remove-orphans

.PHONY: docker-down
docker-down:
	docker-compose -f docker-compose.yaml down --remove-orphans
