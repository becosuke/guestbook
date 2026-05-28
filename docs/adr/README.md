# Architecture Decision Records

設計上の意思決定を時系列に並べる場所。
「何を採用したか」だけでなく「何を採用しなかったか」「なぜそう判断したか」を残し、
将来の議論が同じ前提から始められるようにする。

[`../design/`](../design/) には **現時点での技術設計の意図** だけを置き、ADR は
ここに分離して保持する。design 配下を読めば今の設計が分かり、判断の経緯を辿りたいときは
ここを見る、という運用。

## 運用

- 1 決定につき 1 ファイル。`NNNNNN-kebab-case-title.md` の 6 桁連番形式（`api/configurations/database/migrations/` と同じ桁数）
- フォーマットは Michael Nygard の ADR テンプレート派生:
  - **Status** — proposed / accepted / superseded / deprecated と決定日
  - **Context** — なぜこの判断が必要になったか
  - **Decision** — 採用した内容（断定形）
  - **Alternatives considered** — 採用しなかった選択肢と、それを退けた理由
  - **Consequences** — 採用した結果として受け入れる帰結
- 決定が覆ったときは新しい ADR を起こし、旧 ADR の Status を `superseded by NNNN` に書き換える

## 一覧

| 番号 | タイトル                                                    | Status   |
| ---- | ----------------------------------------------------------- | -------- |
| [000001](./000001-update-post-error-mapping.md) | UpdatePost の状態別エラーコード割り当て | Accepted |
