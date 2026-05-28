# Data Model

## Post

ゲストブックに残された 1 件の投稿を表す唯一のリソース。

### フィールド

| フィールド        | 型                          | 振る舞い | 説明                                                                 |
| ----------------- | --------------------------- | -------- | -------------------------------------------------------------------- |
| `post_id`         | `string` (UUID v4)          | 識別子   | 投稿の一意な ID。サーバ採番。                                         |
| `body`            | `string` (1〜128 文字)      | 入力     | 投稿の本文。前後の空白は正規化されて保存される。                       |
| `valid`           | `bool`                      | 出力専用 | `false` のとき論理削除済み。`false` の場合 `body` は空になる。         |
| `create_time`     | `google.protobuf.Timestamp` | 出力専用 | 作成時刻（サーバ採番）。                                              |
| `update_time`     | `google.protobuf.Timestamp` | 出力専用 | 最終更新時刻（サーバ採番）。新規作成直後は `create_time` と等しい。    |
| `previous_body`   | `string`                    | 出力専用 | 直前 1 世代の本文。新規作成直後は空。論理削除済みでは空。              |

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
              ┌──────────┐
              ▼          │
        ┌──────────┐     │
        │ Created  │ ←───┘
        └────┬─────┘
             │ UpdatePost
             ▼
        ┌──────────┐
        │ Updated  │ ←── UpdatePost（複数回可）
        └────┬─────┘
             │ DeletePost
             ▼
        ┌──────────┐
        │ Deleted  │ ←── 終端状態（以降の Get/Update は失敗または論理削除を示す応答）
        └──────────┘
```

- `Created` / `Updated` の区別は `create_time == update_time` で判定できる（プロダクト的には同一状態として扱う）
- `Deleted` 状態の投稿に対する `UpdatePost` は **拒否**（`FailedPrecondition`）
- `Deleted` 状態の投稿に対する `DeletePost` は **拒否**（`NotFound`）
- `Deleted` 状態の投稿に対する `GetPost` は `valid=false` の `Post` を返す
- `Deleted` 状態の投稿は `ListPosts` の結果には含まれる（`valid=false` として）

### previous_body のセマンティクス

- `Updated` 状態の投稿のみ `previous_body` に値が入る
- 「直前 1 世代」のみを保持し、それより古い本文は失われる
- AIP-148 が示す「変更前後の値を返す」要件を、無制限な履歴を持たずに最小限満たすための割り切り

### 永続化スキーマ

| テーブル        | カラム                                              | 用途                                       |
| --------------- | --------------------------------------------------- | ------------------------------------------ |
| `Posts`         | `PostId` (UUID PK), `PostBody`, `PreviousBody`, `CreateTime`, `UpdateTime`, `DeleteTime` | 投稿本体。`DeleteTime` のゼロ値が「未削除」を示す。 |
| `Paginations`   | `PaginationId` (UUID PK), `Cursor` (JSONB), `CreateTime` | ページネーション状態。`Cursor` は最後に返した行の `(create_time, post_id)` を JSON で保持。 |
