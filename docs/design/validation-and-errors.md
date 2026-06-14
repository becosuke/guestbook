# Validation & Errors

## バリデーションの担当境界

入力の形式チェックは **presentation 層の protovalidate インターセプタ** が一手に担う。
ドメイン層・usecase 層は同じバリデーションを重複して持たない。

### protovalidate で検証されるルール

`proto/guestbook.proto` の `buf.validate` アノテーションで宣言された制約が、
リクエスト受信時にサーバ側のインターセプタで一括検証される。

| リクエスト         | フィールド          | ルール                                      |
| ------------------ | ------------------- | ------------------------------------------- |
| `GetPostRequest`   | `post_id`           | UUID 形式（空文字も拒否）                    |
| `CreatePostRequest`| `post`              | 必須                                        |
|                    | `idempotency_key`   | UUID 形式（空文字も拒否）                    |
| `UpdatePostRequest`| `post`              | 必須                                        |
|                    | `post.post_id`      | 空文字でないこと（メッセージレベル CEL）     |
|                    | `idempotency_key`   | UUID 形式（空文字も拒否）                    |
| `DeletePostRequest`| `post_id`           | UUID 形式（空文字も拒否）                    |
|                    | `idempotency_key`   | UUID 形式（空文字も拒否）                    |
| `ListPostsRequest` | `page_size`         | 0 以上（0 は usecase 層でサーバデフォルト 10 に置換） |
| `Post`(共通)       | `post_id`           | 入力時はゼロ値なら無視、それ以外は UUID 形式 |
|                    | `body`              | 1〜128 文字                                 |

`idempotency_key` には `IGNORE_IF_ZERO_VALUE` を付けていないため、
空文字も `string.uuid` 違反として `INVALID_ARGUMENT` で拒否される。
現状は値そのものを冪等処理に使っていないが、proto 上は常に有効な UUID 文字列を要求する契約になっている。

### handler 側で追加検証されるルール

- `UpdatePostRequest.update_mask` のパスは `body` のみ受理する。それ以外のパスが含まれていれば `INVALID_ARGUMENT`。

## エラーモデル

ドメイン層は sentinel error（`ErrNotFound` 等）を用いて失敗を表現し、presentation 層が gRPC の `status` に詰め替える。
クライアントから見たマッピングは以下のとおり（AIP-193 に準拠）。

| ドメインエラー / 状況              | gRPC code            | REST status |
| ---------------------------------- | -------------------- | ----------- |
| protovalidate のルール違反         | `INVALID_ARGUMENT`   | 400         |
| `update_mask` の不正パス           | `INVALID_ARGUMENT`   | 400         |
| `ErrInvalidArgument`               | `INVALID_ARGUMENT`   | 400         |
| `ErrNotFound`                      | `NOT_FOUND`          | 404         |
| `ErrFailedPrecondition`            | `FAILED_PRECONDITION` | 400         |
| `ErrAlreadyExists`                 | `ALREADY_EXISTS` *1  | 409         |
| `ErrInvalidData`                   | `INTERNAL`           | 500         |
| 上記以外                           | `UNKNOWN`            | 500         |

`ErrInvalidArgument` は「protovalidate を通過する形式上は妥当だが、ドメインルール上は受け付けられない値」が来たときに使う。
状態遷移の前提条件違反は `ErrFailedPrecondition`、外部 IO 越しに観測したデータ破損は `ErrInvalidData` と使い分ける。

*1: `ErrAlreadyExists` はリポジトリ層が PostgreSQL の unique violation を検知した際に返すが、
現状の `CreatePost` ではサーバ採番のため通常のフローでは発生しない（保険的なマッピング）。

### UpdatePost における Precondition / NotFound の使い分け

`UpdatePost` の SQL は WHERE 句に `CreateTime = UpdateTime` を含むため、初回更新後のレコードには
一致しない。0 行更新時はレコードが存在するかで分岐し、振り分けは次のとおり。

| 投稿の状態                                  | 返るエラー              |
| ------------------------------------------- | ----------------------- |
| 未削除で **まだ書き直されていない**         | （成功）                |
| 未削除で **既に 1 度書き直し済み**          | `ErrFailedPrecondition` |
| **論理削除済み**                            | `ErrFailedPrecondition` |
| レコード自体が存在しない                    | `ErrNotFound`           |

判断根拠と却下した代替案（`NOT_FOUND` / `INVALID_ARGUMENT`、ErrorDetails の扱い）は
[../adr/20260529-update-post-error-mapping.md](../adr/20260529-update-post-error-mapping.md) に
ADR として記録している。

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
- `idempotency_key` は値の検証（UUID 形式・必須）こそ行うが、**サーバ側で重複検知・抑止は行わない**。
  受け口の契約だけ整っており、冪等性そのものは保証しない。
