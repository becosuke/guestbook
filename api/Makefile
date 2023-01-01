PROJECT_REPOSITORY=github.com/becosuke/guestbook/api
PROJECT_NAME=guestbook
COMMAND_DIR=internal/cmd
BUILD_DIR=build

CODE_DIRS ?= $(shell go list ./internal/... | grep -v pkg)
TEST_DIRS ?= $(shell go list ./tests/...)
LINT_DIRS ?= $(shell go list ./internal/... | grep -v pkg)

TIMESTAMP := $(shell date "+%Y%m%d%H%M%S")
GIT_REF := $(shell git rev-parse --short=7 HEAD)
VERSION ?= $(TIMESTAMP)-$(GIT_REF)

all: mod test build

mod:
	go mod vendor

mod-tidy:
	go mod tidy

vet:
	go vet $(CODE_DIRS)

lint:
	golangci-lint run $(args) ./internal/...
	go-consistent $(cons_args) $(LINT_DIRS)

fmt:
	@find ./internal -iname "*.go" -not -path "./pkg/**" | xargs gofmt -w

imports:
	@find ./internal -iname "*.go" -not -path "./pkg/**" | xargs goimports -w -local $(PROJECT_REPOSITORY)

run: run-grpc run-rest

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
	rm -f $(BUILD_DIR)/grpc $(BUILD_DIR)/rest

clean-test:
	go clean -testcache

.PHONY: build
build: update-build-dependencies build-grpc build-rest

update-build-dependencies:
	docker run --rm -v $(shell pwd):/root golang:1.18.0-alpine sh -c 'cp /usr/local/go/lib/time/zoneinfo.zip /root/build/zoneinfo.zip && cp /etc/ssl/certs/ca-certificates.crt /root/build/ca-certificates.crt'

.PHONY: build-grpc
build-grpc: build-grpc-linux-binary build-grpc-container

.PHONY: build-rest
build-rest: build-rest-linux-binary build-rest-container

build-grpc-binary:
	CGO_ENABLED=0 go build -a -installsuffix cgo \
	-ldflags "-w -X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-grpc" \
	-o $(BUILD_DIR)/$(PROJECT_NAME)-grpc $(COMMAND_DIR)/grpc/main.go

build-rest-binary:
	CGO_ENABLED=0 go build -a -installsuffix cgo \
	-ldflags "-w -X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-rest" \
	-o $(BUILD_DIR)/$(PROJECT_NAME)-rest $(COMMAND_DIR)/rest/main.go

build-grpc-linux-binary:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
	-ldflags "-w -X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-grpc" \
	-o $(BUILD_DIR)/$(PROJECT_NAME)-grpc $(COMMAND_DIR)/grpc/main.go

build-rest-linux-binary:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
	-ldflags "-w -X main.version=$(VERSION) -X main.serviceName=$(PROJECT_NAME)-rest" \
	-o $(BUILD_DIR)/$(PROJECT_NAME)-rest $(COMMAND_DIR)/rest/main.go

build-grpc-container:
	cd $(BUILD_DIR) && \
	docker build . --no-cache --build-arg name=$(PROJECT_NAME)-grpc -t $(PROJECT_NAME)-grpc:latest -t $(PROJECT_NAME)-grpc:$(VERSION)

build-rest-container:
	cd $(BUILD_DIR) && \
	docker build . --no-cache --build-arg name=$(PROJECT_NAME)-rest -t $(PROJECT_NAME)-rest:latest -t $(PROJECT_NAME)-rest:$(VERSION)

protoc: tools-install
	protoc -I proto -I $(shell brew --prefix)/opt/protobuf/include \
	-I modules/github.com/googleapis/googleapis \
	-I modules/github.com/bufbuild/protoc-gen-validate \
	--go_out pb --go_opt paths=source_relative \
	--go-grpc_out pb --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pb --grpc-gateway_opt paths=source_relative \
	--validate_out "lang=go,paths=source_relative:pb" \
	proto/guestbook.proto

tools-install: tools-tidy
	@for tool in $$(sed -n 's/[ \f\n\r\t]*_ "\(.*\)"/\1/p' tools/tools.go); do go install $${tool}@latest; done

tools-tidy:
	@cd tools && go mod tidy
