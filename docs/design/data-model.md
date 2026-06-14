# Data Model

## Post

ゲストブックに残された 1 件の投稿を表す唯一のリソース。

### フィールド

| フィールド        | 型                          | 振る舞い | 説明                                                                 |
| ----------------- | --------------------------- | -------- | -------------------------------------------------------------------- |
| `post_id`         | `string` (UUID v4)          | 識別子   | 投稿の一意な ID。サーバ採番。                                         |
| `body`            | `string` (1〜128 文字)      | 入力     | 投稿の本文。前後の空白は正規化されて保存される。                       |
| `valid`           | `bool`                      | 出力専用 | `false` のとき論理削除済み。論理削除されている場合は `body` / `create_time` / `update_time` / `previous_body` のいずれも空（未設定）として返す。 |
| `create_time`     | `google.protobuf.Timestamp` | 出力専用 | 作成時刻（サーバ採番）。                                              |
| `update_time`     | `google.protobuf.Timestamp` | 出力専用 | 最終更新時刻（サーバ採番）。書き直し前は `create_time` と等しい。      |
| `previous_body`   | `string`                    | 出力専用 | 書き直し前の本文。新規作成直後は空。論理削除済みでは空。               |

`valid` / `create_time` / `update_time` / `previous_body` は `OUTPUT_ONLY`（AIP-203）であり、
クライアントは作成・更新リクエストに値を含めても無視される（または拒否される）。

### 制約

- `post_id` はサーバ側で UUID v4 を採番する（クライアントは Create リクエストで値を指定しない）
- `body` は 1〜128 文字（protovalidate により presentation 層で検証）
  - 空白のみの本文は `body` の正規化（前後の空白除去）の結果として最小長違反となる
- `body` の前後の空白は正規化（Trim）された値が保存される

### 状態遷移

```
        CreatePost
            │
            ▼
       ┌──────────┐
       │ Created  │ ──── DeletePost ────┐
       └────┬─────┘                     │
            │ UpdatePost                │
            ▼                           │
       ┌──────────┐                     │
       │ Updated  │ ──── DeletePost ───┐│
       └────┬─────┘                    ││
            │ UpdatePost を再度試みても ││
            │ FailedPrecondition で拒否  ││
            ▼                          ▼▼
       ┌──────────┐               ┌──────────┐
       │ Updated  │               │ Deleted  │ ←── 終端状態
       │ (no-op)  │               └──────────┘
       └──────────┘
```

- `Created` / `Updated` の区別は `create_time == update_time` で判定できる
- **UpdatePost は 1 投稿につき 1 回しか成功しない**。`Updated` 状態の投稿に再度 `UpdatePost` を呼ぶと `FailedPrecondition`
- `Deleted` 状態の投稿に対する `UpdatePost` は **拒否**（`FailedPrecondition`）。`Updated` 状態と同じエラーコードに集約されており、両者の区別は現状クライアントに通知していない（[validation-and-errors.md](./validation-and-errors.md) と [../adr/20260529-update-post-error-mapping.md](../adr/20260529-update-post-error-mapping.md) 参照）
- `Deleted` 状態の投稿に対する `DeletePost` は **拒否**（`NotFound`）
- `Deleted` 状態の投稿に対する `GetPost` は `valid=false` の `Post` を返す（`body` 等の他フィールドはすべて空）
- `Deleted` 状態の投稿は `ListPosts` の結果には含まれる（`valid=false` として）
- いずれの操作も `post_id` 自体に該当レコードが無い場合は `NotFound`

### previous_body のセマンティクス

- `Updated` 状態の投稿のみ `previous_body` に値が入る
- 1 投稿につき書き直しは 1 回までのため、`previous_body` は **最初に書かれた本文（初期投稿の本文）** を保持することになる
- AIP-148 が示す「変更前後の値を返す」要件を満たすための保持。書き直しが 1 回までという制約と組み合わさり、結果として「初期投稿 + いまの本文」の 2 点が常に揃う

### 永続化スキーマ

| テーブル        | カラム                                              | 用途                                       |
| --------------- | --------------------------------------------------- | ------------------------------------------ |
| `Posts`         | `PostId` (UUID PK), `PostBody`, `PreviousBody`, `CreateTime`, `UpdateTime`, `DeleteTime` | 投稿本体。`DeleteTime` のゼロ値が「未削除」を示す。 |
| `Paginations`   | `PaginationId` (UUID PK), `Cursor` (JSONB), `CreateTime` | ページネーション状態。`Cursor` は最後に返した行の `(create_time, post_id)` を JSON で保持。 |
