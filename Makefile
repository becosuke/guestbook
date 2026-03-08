GO_VERSION=1.25.7
GO_BINARY=go

.PHONY: build
build: build/api build/view

.PHONY: build-local
build-local: build/api build/view-local

.PHONY: build/api
build/api:
	@$(MAKE) --no-print-directory -C api build

.PHONY: build/view
build/view:
	cd view && docker build -f build/Dockerfile --no-cache -t guestbook-view:latest .

.PHONY: build/view-local
build/view-local:
	cd view && docker build -f build/Dockerfile -t guestbook-view:latest .

.PHONY: buf-generate
buf-generate: buf-dep-update
	buf generate

.PHONY: buf-dep-update
buf-dep-update:
	buf dep update

.PHONY: reset-db
reset-db:
	docker compose rm -sfv postgres
	docker volume rm -f $$(docker volume ls -q --filter name=guestbook_postgres-data)
	docker compose up -d postgres

.PHONY: schema-dump
schema-dump:
	docker compose exec postgres pg_dump -U guestbook -d guestbook --schema-only --no-owner --no-privileges --no-comments -t Posts -t Paginations > api/configurations/database/schema.sql

.PHONY: example/post
example/post:
	curl -v -X POST -H 'Content-Type: application/json' -d '{"post": {"body": "example"}}' 'http://localhost:50080/api/v1/post'

.PHONY: example/put
example/put:
	curl -v -X PUT -H 'Content-Type: application/json' -d '{"post": {"body": "example-value"}}' 'http://localhost:50080/api/v1/post/100'

.PHONY: example/get
example/get:
	curl -v 'http://localhost:50080/api/v1/post/100'

.PHONY: example/list
example/list:
	curl -v 'http://localhost:50080/api/v1/posts/list/10/'

.PHONY: example/delete
example/delete:
	curl -v -X DELETE 'http://localhost:50080/api/v1/post/100'
