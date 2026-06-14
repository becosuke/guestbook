# guestbook

ゲストブック投稿の CRUD API サーバとそのフロントエンドを含むモノレポ。
ポートフォリオとして「どういう原則で設計しているか」を読み取りやすい形で公開することを目的としている。

## Stack

API 定義から UI 動作確認までを 1 リポジトリで完結させ、`make` 1 発でビルドと起動まで終わり、ブラウザから挙動を確認できる構成。

- **API 定義** — Protocol Buffers (`proto/guestbook.proto`)。リクエスト/レスポンスの形に加え、protovalidate によるバリデーションルールも同じファイルで宣言する
- **gRPC 実装** — Go。`buf` で生成した stub の上に `api/internal/cmd/grpc` の Driver / Adapter / Application / Domain を載せる
- **REST 提供** — [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) によるリバースプロキシで、同じ proto から REST/JSON エンドポイントを自動生成。`api/internal/cmd/rest` が起動する
- **Frontend** — TypeScript の簡易 SPA。REST 経由でゲストブックを CRUD するデモ実装で、proto → REST → UI のフロー全体を目視確認できる
- **永続化** — PostgreSQL（compose 同居）
- **一括起動** — `compose.yaml` で gRPC / REST / Frontend / PostgreSQL の 4 サービスをまとめて立ち上げ、`make` だけでビルドから起動まで完了して protobuf 定義から UI 操作までを再現可能

## Architecture

サーバ (`api/`) は Clean Architecture の 4 層構成で、依存は外側から内側への一方通行を厳格に保っている。

```
Driver (cmd/grpc, cmd/rest)
  → Adapter (presentation, repository, infrastructure)
    → Application (usecase)
      → Domain (entities, value objects, errors, interfaces)
```

各層の責務と設計判断の根拠（なぜこの分け方なのか・なぜこの型なのか）はソースコードに隣接した GoDoc コメントとして残してある。ドキュメントを別ファイルに分離せず、実装と一緒に読めるようにする方針。

## Design Principles

### Google AIP に従う

リソース命名、標準メソッド、ページネーション、エラーの扱いは [Google API Improvement Proposals (AIP)](https://google.aip.dev/) を一次根拠にしている。

- **AIP-131〜135** — 標準メソッド (Get / List / Create / Update / Delete) のシグネチャと挙動
- **AIP-134** — `update_mask` (FieldMask) による部分更新
- **AIP-148** — 共通フィールド (`create_time`, `update_time`, `previous_body` 等)
- **AIP-158** — `page_size` / `page_token` を opaque に扱うページネーション
- **AIP-193** — ドメインエラーから gRPC status code へのマッピング

「なぜその型・その振る舞いなのか」を AIP の番号で根拠付けできる状態を維持することで、API デザインの議論を主観に依存させない。

### AI 支援開発を意識した一貫性 / 一意性

LLM によるコード生成と既存コードの整合を取りやすくするため、以下の工夫を入れている。

- **値オブジェクトは typed wrapper** — `PostID` / `PostBody` / `PaginationID` などを薄いラッパ型で表現し、生 string を直接引き回さない
- **コンストラクタを一本化** — `NewPostID` / `NewPost` のような名前のファクトリを 1 つだけ定義し、`PostID(uuid.New())` のような直接 cast を排除する
- **シナリオ別ファクトリ** — `domain.CreatePost` のように利用文脈ごとに名前を分け、`time.Time{}` などゼロ値の引数羅列が呼び出し側に漏れないようにする
- **層を跨ぐ命名の重複を避ける** — `Repositories` / `Querier` / `Commander` 等、grep / シンボル検索したときに意味が一意に取れる名前を選ぶ
- **コメントは WHY だけを残す** — 関数名で語れる WHAT は書かず、AIP との対応・代替案を退けた理由・コンパイラで防げない取り違えの限界など、コードからは読み取りにくい根拠を書く

これにより、AI が既存コードを参照して新しい変更を書くとき、命名・規律・根拠が一致した状態を維持しやすい構造になっている。

## AI 支援開発のセットアップ

このリポジトリは [Claude Code](https://claude.com/claude-code) での開発を前提に、`CLAUDE.md` と `.claude/settings.json` を同梱している。コードベースの概要・ビルド/テストコマンド・アーキテクチャ方針は `CLAUDE.md` から読み取れる。

一方で `.claude/rules/` 配下に置いている核心的な実装ルール（命名・レビュー観点・コミット規約など、判断を縛る詳細）は別の private リポジトリで管理しており、本リポジトリには含めていない（`.gitignore` で除外、ローカルではシンボリックリンクで参照する運用）。

ポートフォリオとして公開しているのは「設計の意図と全体構造」までで、日々の実装を駆動する細則は公開対象外、という線引きにしている。

## Repository structure

- `proto/` — Protocol Buffers スキーマ (`guestbook.proto`)
- `api/` — Go 製の gRPC / gRPC-Gateway サーバ（独立した Go モジュール）
- `view/` — TypeScript フロントエンド
- `third_party/` — git submodules (`googleapis`, `protovalidate`)
- `tools/` — protoc プラグインのバージョン管理
- `aqua.yaml` — CLI ツールのバージョン管理（buf, protoc-gen-go, golangci-lint 等）

## Requirements

- docker
- direnv
- protobuf

## Common commands

リポジトリ直下の Makefile でコード生成・イメージビルド・サンプル呼び出しまで一通り扱える。

```bash
make                # build + up（デフォルトターゲット）
make build          # api (gRPC / REST) と view (frontend) の Docker イメージをビルド
make up             # gRPC / REST / Frontend / PostgreSQL を一括起動
make logs           # 起動中サービスのログを追尾
make down           # 一括停止
make buf-generate   # protobuf コード生成
make schema-dump    # PostgreSQL スキーマを schema.sql にダンプ
make example/post   # サンプル: 投稿の作成
make example/get    # サンプル: 投稿の取得
```

API サーバ単体の開発は `api/Makefile` を使う。

```bash
cd api
make run/grpc       # gRPC サーバ起動
make run/rest       # REST gateway 起動
make test           # ユニットテスト + E2E
make lint           # golangci-lint + go-consistent
make gen            # モック生成 (moq)
```
