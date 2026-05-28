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
| `idempotency_key`   | string | 冪等キー。**必須・UUID 形式**（空文字も拒否される）。  |

### レスポンス

- 成功: 採番済み `post_id`、`create_time = update_time`、`previous_body=""`、`valid=true` の `Post`

### 注意

- `idempotency_key` は **受理されるだけで冪等性を保証する処理には使われていない**。
  サーバ側で重複検知も抑止も行わないため、同一クライアントが連続して送ると毎回新しい `post_id` の投稿が作成される。
- それでも proto バリデーション上は **空文字を含めて省略不可** なので、クライアントは毎回ダミーでも UUID を生成して詰める必要がある。

## UpdatePost

既存 `Post` の本文を差し替える。書き直し前の本文を `previous_body` に退避する。
**1 投稿につき書き直しは 1 回まで**。

### リクエスト

| フィールド          | 型                            | 説明                                                                 |
| ------------------- | ----------------------------- | -------------------------------------------------------------------- |
| `post.post_id`      | string                        | 更新対象。必須（メッセージレベル CEL バリデーションで空文字を拒否）。 |
| `post.body`         | string                        | 新しい本文。1〜128 文字。                                             |
| `idempotency_key`   | string                        | 冪等キー。**必須・UUID 形式**（空文字も拒否される）。                  |
| `update_mask`       | `google.protobuf.FieldMask`   | 更新対象パス。`body` のみ受理。空または未指定でも可。                   |

### update_mask の受理パス

- `body` のみ
- 上記以外のパスを含む場合 `INVALID_ARGUMENT`

### レスポンス

- 成功: 更新後の `Post`（`previous_body` に旧本文、`update_time` が更新される）
- 対象が **既に 1 度書き直し済み**: `FAILED_PRECONDITION`
- 対象が **論理削除済み、もしくは存在しない**: `NOT_FOUND`

### 制約

- `repository/post_commander.go` の UPDATE 文は `CreateTime = UpdateTime` を WHERE 条件に含む。
  初回更新で `UpdateTime` だけが更新されると等値が崩れ、以降の `UpdatePost` は更新行 0 件となり拒否される。
- 結果として `previous_body` に入るのは「最初に書かれた本文」固定になり、それ以降の世代は発生しない。
- `update_mask` を空にしても `body` の更新は実施される。

## DeletePost

論理削除を行う。物理削除はサポートしない。

### リクエスト

| フィールド          | 型     | 説明                                                |
| ------------------- | ------ | --------------------------------------------------- |
| `post_id`           | string | 削除対象の UUID。必須。                              |
| `idempotency_key`   | string | 冪等キー。**必須・UUID 形式**（空文字も拒否される）。 |

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
| `page_size`   | int32  | 1 ページの希望件数。**0 以上**。`0`（または未指定）はサーバデフォルト（10 件）に置き換えられる。負数は `INVALID_ARGUMENT`。 |
| `page_token`  | string | 次ページのトークン。空文字なら最初のページ。Opaque な値。                         |

### レスポンス

| フィールド         | 型           | 説明                                              |
| ------------------ | ------------ | ------------------------------------------------- |
| `posts`            | repeated Post | 結果集合（論理削除済みも含む）。                  |
| `next_page_token`  | string        | 次ページのトークン。空文字なら終端。              |
