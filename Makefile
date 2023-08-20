GO_VERSION=1.20.7
GO_BINARY=go

.PHONY: build
build: build-api

build-api:
	@$(MAKE) --no-print-directory -C api build

.PHONY: protoc
protoc: tools-install
	PATH=$(shell pwd)/bin:$(PATH) protoc -I proto -I $(shell brew --prefix)/opt/protobuf/include \
	-I modules/github.com/googleapis/googleapis \
	-I modules/github.com/bufbuild/protoc-gen-validate \
	--go_out pbgo --go_opt paths=source_relative \
	--go-grpc_out pbgo --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pbgo --grpc-gateway_opt paths=source_relative \
	--validate_out "lang=go,paths=source_relative:pbgo" \
	--openapiv2_out . \
	proto/guestbook.proto

tools-install: tools-tidy
	@for tool in $$(sed -n 's/[ \f\n\r\t]*_ "\(.*\)"/\1/p' tools/tools.go); do GOBIN=$(shell pwd)/bin $(GO_BINARY) install $${tool}@latest; done

tools-tidy:
	@cd tools && $(GO_BINARY) mod tidy

pbgo-tidy:
	@cd pbgo && $(GO_BINARY) mod tidy

example-post:
	curl -v -X POST -H 'Content-Type: application/json' -d '{"post": {"body": "example"}}' 'http://localhost:50080/api/v1/post'

example-put:
	curl -v -X PUT -H 'Content-Type: application/json' -d '{"post": {"body":"example-value"}}' 'http://localhost:50080/api/v1/post/100'

example-get:
	curl -v 'http://localhost:50080/api/v1/post/100'

example-delete:
	curl -v -X DELETE 'http://localhost:50080/api/v1/post/100'
