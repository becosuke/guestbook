# Configuration

サービスの起動時設定と、デプロイにあたり外から触れる表面の一覧。

## 環境変数

すべて 12-factor 的に環境変数経由で注入する。値は `envconfig` で解決される。

| 変数             | 既定値            | 用途                                              |
| ---------------- | ----------------- | ------------------------------------------------- |
| `ENVIRONMENT`    | `development`     | 実行環境。`development` / `production` / `test`。 |
| `LOG_LEVEL`      | `debug`           | zap のログレベル（`debug`/`info`/`warn`/`error`）。 |
| `GRPC_HOST`      | `127.0.0.1`       | gRPC サーバの bind ホスト。                        |
| `GRPC_PORT`      | `50051`           | gRPC サーバの bind ポート。                        |
| `REST_HOST`      | `127.0.0.1`       | REST ゲートウェイの bind ホスト。                  |
| `REST_PORT`      | `50080`           | REST ゲートウェイの bind ポート。                  |
| `DATABASE_URL`   | (なし)            | PostgreSQL 接続文字列（必須）。                    |

### 値の取り扱い

- `ENVIRONMENT` は大文字小文字を吸収する（`Development` / `PRODUCTION` も受理）。
  未知の値が来た場合は `development` にフォールバックする。
- `LOG_LEVEL` を解釈できない値が来た場合は `debug` にフォールバックする。
- `DATABASE_URL` 未指定でも起動プロセス自体は走るが、DB アクセス時にエラーになる
  （`grpc/main.go` 起動時のプール作成で fatal 終了）。

## プロセス構成

API は **gRPC プロセス** と **REST ゲートウェイプロセス** の 2 つに分かれている。

```
[ Client ] --HTTP/JSON--> [ rest プロセス :50080 ]
                                 │
                                 │ gRPC
                                 ▼
              [ grpc プロセス :50051 ] --SQL--> [ PostgreSQL ]
```

- REST プロセスはアプリケーションロジックを持たず、gRPC への薄いプロキシ。
- gRPC プロセスのみがデータベースに接続する。
- 認証・バリデーションは gRPC プロセスのインターセプタで実施される。

## Docker Compose 構成

`compose.yaml` は 4 サービスを起動する。

| サービス   | イメージ              | 公開ポート | 役割                                  |
| ---------- | --------------------- | ---------- | ------------------------------------- |
| `grpc`     | `guestbook-api-grpc`  | （内部）   | gRPC API サーバ                         |
| `rest`     | `guestbook-api-rest`  | `50080`    | gRPC-Gateway（REST → gRPC ブリッジ）   |
| `view`     | `guestbook-view`      | `8000`     | デモ用フロントエンド（REST を叩く SPA） |
| `postgres` | `postgres:17-alpine`  | `5432`     | データストア                            |

PostgreSQL は初期化時に `api/configurations/database/schema.sql` を読み込み、
`Posts` / `Paginations` テーブルを作成する。マイグレーションは
`api/configurations/database/migrations/` 配下に SQL ファイル形式で管理されている。

## 観測性

- 構造化ログ: zap（`grpc_zap` インターセプタが gRPC リクエスト単位のログを出す）
- メトリクス、トレース、ヘルスチェックエンドポイントは現状提供しない
- PostgreSQL のヘルスチェックは Compose の `pg_isready` で確認している
