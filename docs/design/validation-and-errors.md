# Validation & Errors

## バリデーションの担当境界

入力の形式チェックは **presentation 層の protovalidate インターセプタ** が一手に担う。
ドメイン層・usecase 層は同じバリデーションを重複して持たない。

### protovalidate で検証されるルール

`proto/guestbook.proto` の `buf.validate` アノテーションで宣言された制約が、
リクエスト受信時にサーバ側のインターセプタで一括検証される。

| リクエスト         | フィールド          | ルール                                      |
| ------------------ | ------------------- | ------------------------------------------- |
| `GetPostRequest`   | `post_id`           | UUID 形式                                   |
| `CreatePostRequest`| `post`              | 必須                                        |
|                    | `idempotency_key`   | UUID 形式                                   |
| `UpdatePostRequest`| `post`              | 必須                                        |
|                    | `post.post_id`      | 空文字でないこと（メッセージレベル CEL）     |
|                    | `idempotency_key`   | UUID 形式                                   |
| `DeletePostRequest`| `post_id`           | UUID 形式                                   |
|                    | `idempotency_key`   | UUID 形式                                   |
| `ListPostsRequest` | `page_size`         | 0 より大きい                                |
| `Post`(共通)       | `post_id`           | 入力時はゼロ値なら無視、それ以外は UUID 形式 |
|                    | `body`              | 1〜128 文字                                 |

### handler 側で追加検証されるルール

- `UpdatePostRequest.update_mask` のパスは `body` のみ受理する。それ以外のパスが含まれていれば `INVALID_ARGUMENT`。

## エラーモデル

ドメイン層は sentinel error（`ErrNotFound` 等）を用いて失敗を表現し、presentation 層が gRPC の `status` に詰め替える。
クライアントから見たマッピングは以下のとおり（AIP-193 に準拠）。

| ドメインエラー / 状況              | gRPC code            | REST status |
| ---------------------------------- | -------------------- | ----------- |
| protovalidate のルール違反         | `INVALID_ARGUMENT`   | 400         |
| `update_mask` の不正パス           | `INVALID_ARGUMENT`   | 400         |
| `ErrNotFound`                      | `NOT_FOUND`          | 404         |
| `ErrFailedPrecondition`            | `FAILED_PRECONDITION` | 400         |
| `ErrAlreadyExists`                 | `ALREADY_EXISTS` *1  | 409         |
| `ErrInvalidArgument` / `ErrInvalidData` | `INTERNAL`     | 500         |
| 上記以外                           | `UNKNOWN`            | 500         |

*1: `ErrAlreadyExists` はリポジトリ層が PostgreSQL の unique violation を検知した際に返すが、
現状の `CreatePost` ではサーバ採番のため通常のフローでは発生しない（保険的なマッピング）。

## エラー詳細（`google.rpc.BadRequest`）

protovalidate のバリデーション違反では、`status.details` に `google.rpc.BadRequest` を埋め込む。
クライアントは違反フィールドごとの理由を読み取れる。

```
status:
  code: INVALID_ARGUMENT
  message: "validation error"
  details:
    - "@type": type.googleapis.com/google.rpc.BadRequest
      field_violations:
        - field: "post.body"
          description: "value length must be at least 1 characters"
```

## 操作ごとの可能なエラー

| 操作          | INVALID_ARGUMENT | NOT_FOUND | FAILED_PRECONDITION | INTERNAL |
| ------------- | :--------------: | :-------: | :-----------------: | :------: |
| `GetPost`     | ✓                | ✓         |                     | ✓        |
| `CreatePost`  | ✓                |           |                     | ✓        |
| `UpdatePost`  | ✓                | ✓         | ✓                   | ✓        |
| `DeletePost`  | ✓                | ✓         |                     | ✓        |
| `ListPosts`   | ✓                | ✓         |                     | ✓        |

## 現状の制約事項

- 認証は **インターセプタが定義されているのみで実質ノーオペレーション**（`authFunc` は context をそのまま返す）。
  認可・テナント分離は提供されない。
- `idempotency_key` は **受理のみ**。サーバ側で重複検知・抑止は行わない。
