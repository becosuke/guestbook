# Guestbook 設計書

このディレクトリは、Guestbook サービスの技術設計を扱う。
「どう作ってあるか・なぜそう作ったか」の根拠を残すことが目的で、
プロダクトとして提供する機能・振る舞いの仕様は [`../spec/`](../spec/) を参照する。

design 配下には **現時点の意図のみ** を置く。
意思決定の履歴・採用と却下の根拠は [`../adr/`](../adr/) を参照する。

## 目次

- [data-model.md](./data-model.md) — リソースモデル `Post` のフィールドと制約
- [api-operations.md](./api-operations.md) — gRPC / REST の標準メソッドと振る舞い
- [pagination.md](./pagination.md) — `ListPosts` のページネーション方式
- [validation-and-errors.md](./validation-and-errors.md) — 入力バリデーションとエラー表現
- [configuration.md](./configuration.md) — 起動時設定とデプロイ構成

## 設計根拠

API スキーマと振る舞いの設計判断は [Google API Improvement Proposals (AIP)](https://google.aip.dev/) を一次根拠としている。
個別の設計書では「どの AIP に従っているか」を明示する。

| 領域             | 一次根拠                                          |
| ---------------- | ------------------------------------------------- |
| 標準メソッド     | AIP-131〜135（Get / List / Create / Update / Delete） |
| 部分更新         | AIP-134（`update_mask` / `FieldMask`）             |
| 共通フィールド   | AIP-148（`create_time`, `update_time` ほか）       |
| ページネーション | AIP-158（`page_size`, `page_token` の opaque 性）   |
| エラーモデル     | AIP-193（ドメインエラー → gRPC status code）        |
