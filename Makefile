GO_VERSION=1.25.7
GO_BINARY=go

.PHONY: build
build: build-api

build-api:
	@$(MAKE) --no-print-directory -C api build

.PHONY: buf-generate
buf-generate: buf-dep-update
	buf generate

buf-dep-update:
	buf dep update

example-post:
	curl -v -X POST -H 'Content-Type: application/json' -d '{"post": {"body": "example"}}' 'http://localhost:50080/api/v1/post'

example-put:
	curl -v -X PUT -H 'Content-Type: application/json' -d '{"post": {"body": "example-value"}}' 'http://localhost:50080/api/v1/post/100'

example-get:
	curl -v 'http://localhost:50080/api/v1/post/100'

example-delete:
	curl -v -X DELETE 'http://localhost:50080/api/v1/post/100'
