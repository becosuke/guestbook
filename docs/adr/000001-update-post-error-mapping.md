# 000001. UpdatePost の状態別エラーコード割り当て

## Status

Accepted (2026-05-29)

## Context

`UpdatePost` は実装上、1 投稿につき 1 回だけ書き直しを認める設計になっている
（`repository/post_commander.go` の UPDATE 文に `CreateTime = UpdateTime` 条件が入っているため）。
また論理削除済みの投稿も同様に更新できない。

その結果、更新が拒否される状況は 3 つに分かれる:

1. 投稿が存在しない（`post_id` 自体が見つからない）
2. 投稿は存在するが、既に 1 度書き直し済み
3. 投稿は存在するが、論理削除済み

1 は「リソース自体が見つからない」、2 と 3 は「リソースは存在するが現在の状態が更新を許さない」と
性質が異なる。どの gRPC code に割り当てるかを決める必要があり、特に 3（削除済み）の扱いで
議論があった。

候補は次の 3 通り:

- `NOT_FOUND` — 論理削除を「実質的に消えた」と扱う立場
- `FAILED_PRECONDITION` — 「リソースは存在し、状態が更新を許さない」と扱う立場
- `INVALID_ARGUMENT` — 「永続的に更新不能なリソースへの参照は引数として不正」と扱う立場

## Decision

3（削除済み）のケースは **`FAILED_PRECONDITION`** を返す。
あわせて、2（書き直し済み）も同じ `FAILED_PRECONDITION` に集約する。
1（存在しない）のみ `NOT_FOUND`。

`UpdatePost` のエラー応答は以下のマッピングになる:

| 状況                            | gRPC code             |
| ------------------------------- | --------------------- |
| 投稿が存在しない                | `NOT_FOUND`           |
| 既に 1 度書き直し済み           | `FAILED_PRECONDITION` |
| 論理削除済み                    | `FAILED_PRECONDITION` |

2 と 3 の区別は **クライアントには通知しない**。区別が必要になった場合は
`google.rpc.PreconditionFailure` などの ErrorDetails で理由を併送する方針とするが、
現時点では ErrorDetails の付与は実装しない。

## Alternatives considered

### `NOT_FOUND` for deleted post

論理削除を「ユーザから見ると hard delete と等価」に扱う設計とは整合する。

**不採用の理由**: 本プロジェクトは `GetPost` が論理削除済みの投稿を `valid=false` の `Post` として
返す設計を採っており、削除済みでも「観測可能なリソース」として扱っている。
Update でだけ「存在しない」体を取ると、`Get` の振る舞いと矛盾する。

### `INVALID_ARGUMENT` for deleted post

「論理削除は戻らないというビジネスルールから、削除済みリソースへの参照は永続的に通らない →
引数として不正」と捉える見方。

**不採用の理由**:

- `post_id` の形式は妥当（UUID）、かつ参照先のリソースは存在する。
  `INVALID_ARGUMENT` は「引数自体が壊れている」を表すコードであり、リソース状態を理由にした
  拒否は本来この枠の外。
- 公式の `FAILED_PRECONDITION` 説明にある「retry until the state has been explicitly fixed」を
  厳格に読むと、戻らない永続状態には合わないという議論はあり得る。ただし Google 自身の API
  （Compute Engine の TERMINATED instance、BigQuery の completed job など）でも永続状態に対して
  `FAILED_PRECONDITION` を使っており、運用上はこちらが標準的な扱い。

### 2 ケースをそれぞれ別の gRPC code に分ける

「書き直し済み」と「削除済み」をクライアントから直接区別したいなら、片方を
`FAILED_PRECONDITION`、もう片方を別コード（例えば `ABORTED`）に分ける選択肢もある。

**不採用の理由**: gRPC 公式コードに「書き直し済み」「削除済み」をそれぞれ自然に表す枠は無い。
状態の細分は code ではなく `google.rpc.PreconditionFailure` などの ErrorDetails で表現するのが
標準的な方法（AIP-193 が示す方向）。

## Consequences

- クライアントは `FAILED_PRECONDITION` が返ったとき、「書き直し済み」か「削除済み」かを区別できない。
  ただし、どちらの場合も「リトライしても通らない」「クライアント側の操作で状態を戻す手段はない」
  という点は共通しており、UX 上の差は限定的。
- 将来この区別が必要になった場合は、ErrorDetails の `PreconditionFailure` を付与して理由を伝える。
  proto / public API の変更は不要で、後方互換的に追加できる。
- `Get` / `Update` の論理削除に対する振る舞いが「Get は `valid=false` で返す、Update は
  `FAILED_PRECONDITION` で拒否する」と一貫し、「削除済みリソースが観測可能」というセマンティクスが
  layer をまたいで揃う。
