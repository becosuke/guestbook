.PHONY: build
build: build-api

build-api:
	@$(MAKE) --no-print-directory -C api build

.PHONY: protoc
protoc: tools-install
	PATH=$(shell pwd)/bin:$(PATH) protoc -I proto -I $(shell brew --prefix)/opt/protobuf/include \
	-I modules/github.com/googleapis/googleapis \
	-I modules/github.com/bufbuild/protoc-gen-validate \
	--go_out pb --go_opt paths=source_relative \
	--go-grpc_out pb --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pb --grpc-gateway_opt paths=source_relative \
	--validate_out "lang=go,paths=source_relative:pb" \
	--openapiv2_out . \
	proto/guestbook.proto

tools-install: tools-tidy
	@for tool in $$(sed -n 's/[ \f\n\r\t]*_ "\(.*\)"/\1/p' tools/tools.go); do GOBIN=$(shell pwd)/bin go install $${tool}@latest; done

tools-tidy:
	@cd tools && go mod tidy
