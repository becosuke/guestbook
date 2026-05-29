# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Guestbook is a Go-based gRPC/REST API server with a clean architecture design. It exposes a CRUD API for "Post" resources via both gRPC and gRPC-Gateway (REST).

## Implementation Rules

核心的な実装ルールは `.claude/rules/` 配下のファイルに置く。作業前に該当ファイルを読むこと。

- `.claude/rules/` は別リポジトリのシンボリックリンクで管理しており、`.gitignore` に登録済み（本リポジトリにはコミットされない）

## Repository Structure

- `proto/` — Protocol Buffers definitions (single `guestbook.proto`)
- `api/` — Go API server (separate Go module: `github.com/becosuke/guestbook/api`)
- `third_party/` — Git submodules (`googleapis`, `protovalidate`)
- `tools/` — Root-level Go tool dependencies for protoc plugins
- `aqua.yaml` — CLI tool version management (buf, protoc-gen-go, golangci-lint, etc.)

## Architecture (api/)

Clean architecture with 4 layers. Dependencies flow inward only.

```
Driver (cmd/grpc, cmd/rest)
  → Adapter (presentation, repository, infrastructure)
    → Application (usecase)
      → Domain (entities, value objects, errors, interfaces)
```

- **Domain** (`internal/domain/`): Entities and value objects (`Post`, `PostID`, `PostBody`, `PageOption`, `Config`, `Environment`), domain errors
  - `interfaces/`: Repository interfaces (`Querier`, `Commander`) and mocks
- **Application** (`internal/usecase/`): `Usecase` struct with business logic
- **Adapter** (`internal/adapter/`):
  - `presentation/`: gRPC handler, `converter.go` for domain ↔ protobuf conversion, `Usecase` interface
  - `repository/`: PostgreSQL implementation of `Querier` and `Commander`
  - `infrastructure/config/`: Environment variable loading via envconfig
- **Driver** (`internal/cmd/`): gRPC server (`grpc/`) and REST gateway (`rest/`) entry points with middleware chain (zap logging, auth, validation, recovery)
- **pkg/** (`internal/pkg/`): Shared packages — `logger` (zap), `pb` (generated protobuf code)

Dependency injection is done manually in `internal/cmd/grpc/main.go` and `internal/cmd/rest/main.go`.

## Build Rules

- `go build` を直接実行しないこと。ビルドは必ず `make` を通して行う（`make build`, `make run/grpc` など）

## Build & Dev Commands

### Root Makefile (from repo root)

```bash
make buf-generate      # Generate all protobuf code via buf
make schema-dump       # Dump PostgreSQL schema to schema.sql
make example/post      # Example: create a post
make example/get       # Example: get a post
make example/put       # Example: update a post
make example/delete    # Example: delete a post
```

### API Makefile (from api/)

```bash
make mod               # Download Go modules
make test              # Run all tests (unit + e2e)
make test/unit         # Unit tests only (internal/...)
make test/e2e          # E2E tests only (test/...)
make run/grpc          # Run gRPC server locally
make run/rest          # Run REST gateway locally
make lint              # golangci-lint + go-consistent
make fmt               # Format with goimports
make vet               # go vet
make gen               # Generate mocks (moq via go generate)
make build             # Build Docker images for grpc + rest
```

### Running a single test

```bash
cd api && go test -run TestName ./internal/adapter/presentation/...
```

## Proto Code Generation

Proto source: `proto/guestbook.proto`. Uses `buf` for code generation (`buf.yaml`, `buf.gen.yaml`):
- `protoc-gen-go` / `protoc-gen-go-grpc` — Go message/service stubs
- `protoc-gen-grpc-gateway` — REST HTTP gateway
- `protoc-gen-openapiv2` — OpenAPI spec
- `protoc-gen-validate` — Validation rules (protovalidate)

Generated Go code output: `api/internal/pkg/pb/`

Proto dependencies are managed via `buf.yaml` (`buf.build/googleapis/googleapis`, `buf.build/bufbuild/protovalidate`).

## Configuration

Environment variables (defaults):
- `ENVIRONMENT` (development), `LOG_LEVEL` (info)
- `GRPC_HOST` (""), `GRPC_PORT` (50051)
- `REST_HOST` (""), `REST_PORT` (50080)
- `DATABASE_URL` (no default)

## Testing Patterns

- `moq` for mocking repository and usecase interfaces
- Mocks generated via `moq`, stored alongside interfaces (`api/internal/domain/interfaces/`, `api/internal/adapter/presentation/`)
- Assertions use `testify` (`assert`, `require`)

## API Design Guidelines

This project follows [Google API Improvement Proposals (AIP)](https://google.aip.dev/). Key conventions:
- Resource names and field names follow AIP naming standards (e.g., `create_time`, `update_time`)
- Standard methods (Create, Get, Update, Delete, List) follow AIP-131 through AIP-135
- Field behavior annotations follow AIP-203

## Tool Management

CLI tools managed by [aqua](https://aquaproj.github.io/) (`aqua.yaml`). Go-based protoc plugins also listed in `tools/tools.go` for `make tools-install`.
