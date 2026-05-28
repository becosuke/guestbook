# API Operations

Guestbook が提供する標準メソッドの一覧と振る舞い。サービス名は `GuestbookService`、
proto 定義は `proto/guestbook.proto` を参照する。

| メソッド       | gRPC                | REST                                                   | AIP     |
| -------------- | ------------------- | ------------------------------------------------------ | ------- |
| `GetPost`      | `GuestbookService/GetPost`      | `GET /api/v1/post/{post_id}`                | AIP-131 |
| `CreatePost`   | `GuestbookService/CreatePost`   | `POST /api/v1/post`                          | AIP-133 |
| `UpdatePost`   | `GuestbookService/UpdatePost`   | `PATCH /api/v1/post/{post.post_id}`          | AIP-134 |
| `DeletePost`   | `GuestbookService/DeletePost`   | `DELETE /api/v1/post/{post_id}`              | AIP-135 |
| `ListPosts`    | `GuestbookService/ListPosts`    | `GET /api/v1/posts/list/{page_size}/{page_token}` | AIP-132 |

REST のルーティングは gRPC-Gateway によって proto 上の `google.api.http` アノテーションから自動生成される。

## GetPost

論理削除済み（`valid=false`）であってもレコードが存在する限り `Post` を返す。
削除済みの場合 `body` は空となる。

### リクエスト

| フィールド | 型     | 説明                              |
| ---------- | ------ | --------------------------------- |
| `post_id`  | string | 取得対象の UUID。必須（UUID 形式）。 |

### レスポンス

- 成功: `Post`
- `post_id` に該当するレコードが存在しない: `NOT_FOUND`

## CreatePost

新規 `Post` を作成し、採番された `post_id` を含む `Post` を返す。

### リクエスト

| フィールド          | 型     | 説明                                                  |
| ------------------- | ------ | ----------------------------------------------------- |
| `post.body`         | string | 本文。1〜128 文字。必須。                              |
| `post.post_id`      | string | クライアントから指定しても無視される（サーバ採番）。   |
| `idempotency_key`   | string | 冪等キー。指定する場合 UUID 形式（現状は受理のみ）。   |

### レスポンス

- 成功: 採番済み `post_id`、`create_time = update_time`、`previous_body=""`、`valid=true` の `Post`

### 注意

- `idempotency_key` は **受理時点では冪等性を保証する処理に使われていない**（受け口のみ提供）
- 同一クライアントが連続して送ると毎回新しい `post_id` の投稿が作成される

## UpdatePost

既存 `Post` の本文を差し替える。**直前 1 世代** の本文を `previous_body` に退避する。

### リクエスト

| フィールド          | 型                            | 説明                                                                 |
| ------------------- | ----------------------------- | -------------------------------------------------------------------- |
| `post.post_id`      | string                        | 更新対象。必須（メッセージレベル CEL バリデーションで空文字を拒否）。 |
| `post.body`         | string                        | 新しい本文。1〜128 文字。                                             |
| `idempotency_key`   | string                        | 冪等キー。指定する場合 UUID 形式（現状は受理のみ）。                   |
| `update_mask`       | `google.protobuf.FieldMask`   | 更新対象パス。`body` のみ受理。空または未指定でも可。                   |

### update_mask の受理パス

- `body` のみ
- 上記以外のパスを含む場合 `INVALID_ARGUMENT`

### レスポンス

- 成功: 更新後の `Post`（`previous_body` に旧本文、`update_time` が更新される）
- 対象が論理削除済み: `FAILED_PRECONDITION`
- 対象が存在しない: `NOT_FOUND`

### 制約

- 「2 世代以上前」の本文は失われる（履歴テーブルを持たない設計）
- `update_mask` を空にしても `body` の更新は実施される

## DeletePost

論理削除を行う。物理削除はサポートしない。

### リクエスト

| フィールド          | 型     | 説明                                                |
| ------------------- | ------ | --------------------------------------------------- |
| `post_id`           | string | 削除対象の UUID。必須。                              |
| `idempotency_key`   | string | 冪等キー。指定する場合 UUID 形式(現状は受理のみ)。   |

### レスポンス

- 成功: `google.protobuf.Empty`
- 既に論理削除済み、または存在しない: `NOT_FOUND`

### 注意

- 論理削除のため、削除後の `GetPost` は `valid=false` のレスポンスを返す
- 削除済み投稿は `ListPosts` の結果集合に **含まれる**（`valid=false` として返却される）

## ListPosts

`Post` を作成時刻の降順（同時刻は `post_id` 昇順）で列挙する。
詳細な page token のセマンティクスは [pagination.md](./pagination.md) を参照。

### リクエスト

| フィールド    | 型     | 説明                                                                            |
| ------------- | ------ | ------------------------------------------------------------------------------- |
| `page_size`   | int32  | 1 ページの希望件数。1 以上。未指定（0）の場合はサーバデフォルト（10）を使用。      |
| `page_token`  | string | 次ページのトークン。空文字なら最初のページ。Opaque な値。                         |

### レスポンス

| フィールド         | 型           | 説明                                              |
| ------------------ | ------------ | ------------------------------------------------- |
| `posts`            | repeated Post | 結果集合（論理削除済みも含む）。                  |
| `next_page_token`  | string        | 次ページのトークン。空文字なら終端。              |
