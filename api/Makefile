PROJECT_REPOSITORY=github.com/becosuke/guestbook/api
PROJECT_NAME=guestbook-api
COMMAND_DIR=internal/cmd
BUILD_DIR=build
GO_VERSION=1.18.10
GO_BINARY=go$(GO_VERSION)

CODE_DIRS ?= $(shell $(GO_BINARY) list ./internal/... | grep -v pkg)
TEST_DIRS ?= $(shell $(GO_BINARY) list ./tests/...)

TIMESTAMP := $(shell date "+%Y%m%d%H%M%S")
GIT_REF := $(shell git rev-parse --short=7 HEAD)
VERSION ?= $(TIMESTAMP)-$(GIT_REF)

all: mod test build

mod:
	$(GO_BINARY) mod vendor

mod-tidy:
	$(GO_BINARY) mod tidy

vet:
	$(GO_BINARY) vet $(CODE_DIRS)

lint:
	golangci-lint run $(args) ./internal/...
	go-consistent $(cons_args) $(CODE_DIRS)

fmt:
	@find ./internal -iname "*.go" -not -path "./pkg/**" | xargs gofmt -w

imports:
	@find ./internal -iname "*.go" -not -path "./pkg/**" | xargs goimports -w -local $(PROJECT_REPOSITORY)

run: mod
	$(GO_BINARY) run -ldflags "-X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-server" \
	$(PROJECT_REPOSITORY)/$(COMMAND_DIR)/server

.PHONY: test
test: unit-test functional-test

unit-test:
	$(GO_BINARY) test $(args) -race -cover $(CODE_DIRS)

functional-test:
	$(GO_BINARY) test $(args) -race -cover $(TEST_DIRS)

.PHONY: clean
clean:
	$(GO_BINARY) clean
	rm -f $(BUILD_DIR)/server

clean-test:
	$(GO_BINARY) clean -testcache

.PHONY: build
build: update-build-dependencies build-server

update-build-dependencies:
	docker run --rm -v $(shell pwd):/root golang:1.18.10-alpine sh -c 'cp /usr/local/go/lib/time/zoneinfo.zip /root/build/zoneinfo.zip && cp /etc/ssl/certs/ca-certificates.crt /root/build/ca-certificates.crt'

.PHONY: build-server
build-server: build-server-binary build-server-container

build-server-binary:
	CGO_ENABLED=0 GOOS=linux $(GO_BINARY) build -a -installsuffix cgo \
	-ldflags "-w -X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-server" \
	-o $(BUILD_DIR)/$(PROJECT_NAME)-server $(COMMAND_DIR)/server/main.go

build-server-container:
	cd $(BUILD_DIR) && \
	docker build . --no-cache --build-arg name=$(PROJECT_NAME)-server -t $(PROJECT_NAME)-server:latest -t $(PROJECT_NAME)-server:$(VERSION)

gen-mocks:
	mockgen -source internal/domain/post/usecase.go -destination mocks/domain/post/usecase.go -package post
	mockgen -source internal/domain/post/repository.go -destination mocks/domain/post/repository.go -package post
	mockgen -source internal/domain/post/generator.go -destination mocks/domain/post/generator.go -package post
	mockgen -source internal/drivers/syncmap/syncmap.go -destination mocks/drivers/syncmap/syncmap.go -package syncmap

tools-install: tools-tidy
	@for tool in $$(sed -n 's/[ \f\n\r\t]*_ "\(.*\)"/\1/p' tools/tools.go); do GOBIN=$(shell pwd)/bin go install $${tool}@latest; done

tools-tidy:
	@cd tools && $(GO_BINARY) mod tidy
