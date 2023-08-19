GO_VERSION=1.20.7
GO_BINARY=go

PROJECT_REPOSITORY_ORIGIN=github.com/becosuke/guestbook
PROJECT_REPOSITORY=$(PROJECT_REPOSITORY_ORIGIN)/api
PROJECT_NAME=guestbook-api
COMMAND_DIR=internal/cmd
BUILD_DIR=build

CODE_DIRS ?= $(shell $(GO_BINARY) list ./internal/... | grep -v pkg)
TEST_DIRS ?= $(shell $(GO_BINARY) list ./tests/...)

TIMESTAMP := $(shell date "+%Y%m%d%H%M%S")
GIT_REF := $(shell git rev-parse --short=7 HEAD)
VERSION ?= $(TIMESTAMP)-$(GIT_REF)

all: mod test build

mod:
	go mod download

mod-tidy:
	go mod tidy

vet:
	go vet $(CODE_DIRS)

lint:
	./bin/golangci-lint run $(args) ./internal/...
	./bin/go-consistent $(cons_args) $(CODE_DIRS)

fmt:
	@find ./internal -iname "*.go" -not -path "./pkg/**" | xargs ./bin/goimports -w -local $(PROJECT_REPOSITORY_ORIGIN)

run: mod run-grpc run-rest

run-grpc:
	go run -ldflags "-X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-grpc" \
	$(PROJECT_REPOSITORY)/$(COMMAND_DIR)/grpc

run-rest:
	go run -ldflags "-X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-rest" \
	$(PROJECT_REPOSITORY)/$(COMMAND_DIR)/rest

.PHONY: test
test: unit-test functional-test

unit-test:
	go test $(args) -race -cover $(CODE_DIRS)

functional-test:
	go test $(args) -race -cover $(TEST_DIRS)

.PHONY: clean
clean:
	go clean
	rm -f $(BUILD_DIR)/zoneinfo.zip $(BUILD_DIR)/ca-certificates.crt $(BUILD_DIR)/$(PROJECT_NAME)-grpc $(BUILD_DIR)/$(PROJECT_NAME)-rest

clean-test:
	go clean -testcache

.PHONY: build
build: clean update-build-dependencies build-grpc build-rest

update-build-dependencies:
	docker run --rm -v $(shell pwd):/root golang:$(GO_VERSION)-alpine sh -c 'cp /usr/local/go/lib/time/zoneinfo.zip /root/build/zoneinfo.zip && cp /etc/ssl/certs/ca-certificates.crt /root/build/ca-certificates.crt'

build-grpc: build-grpc-binary build-grpc-container

build-grpc-binary:
	CGO_ENABLED=0 GOOS=linux $(GO_BINARY) build -a -installsuffix cgo \
	-ldflags "-w -X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-grpc" \
	-o $(BUILD_DIR)/$(PROJECT_NAME)-grpc $(COMMAND_DIR)/grpc/main.go

build-grpc-container:
	cd $(BUILD_DIR) && \
	docker build . --no-cache --build-arg name=$(PROJECT_NAME)-grpc -t $(PROJECT_NAME)-grpc:latest -t $(PROJECT_NAME)-grpc:$(VERSION)

build-rest: build-rest-binary build-rest-container

build-rest-binary:
	CGO_ENABLED=0 GOOS=linux $(GO_BINARY) build -a -installsuffix cgo \
	-ldflags "-w -X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-rest" \
	-o $(BUILD_DIR)/$(PROJECT_NAME)-rest $(COMMAND_DIR)/rest/main.go

build-rest-container:
	cd $(BUILD_DIR) && \
	docker build . --no-cache --build-arg name=$(PROJECT_NAME)-rest -t $(PROJECT_NAME)-rest:latest -t $(PROJECT_NAME)-rest:$(VERSION)

gen: gen-mocks

gen-mocks:
	mockgen -source internal/domain/post/usecase.go -destination mocks/domain/post/usecase.go -package post
	mockgen -source internal/domain/post/repository.go -destination mocks/domain/post/repository.go -package post
	mockgen -source internal/domain/post/generator.go -destination mocks/domain/post/generator.go -package post
	mockgen -source internal/drivers/syncmap/syncmap.go -destination mocks/drivers/syncmap/syncmap.go -package syncmap

tools-install: tools-tidy
	@for tool in $$(sed -n 's/[ \f\n\r\t]*_ "\(.*\)"/\1/p' tools/tools.go); do GOBIN=$(shell pwd)/bin $(GO_BINARY) install $${tool}@latest; done

tools-tidy:
	@cd tools && $(GO_BINARY) mod tidy
