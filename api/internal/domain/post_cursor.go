package domain

import (
	"encoding/json"
	"time"
)

// PostCursor は「次ページの開始位置」を表すドメイン値。
//
// seek method 方式のページネーションを採用しており、(CreateTime DESC, PostID ASC)
// で一意に並ぶ Posts テーブルに対して「直前に取得した最後の行」をカーソルとして
// 持つ。OFFSET 方式と比べたメリットは、
//   - データ量が増えても性能が劣化しない（O(log n) のインデックスシークで済む）
//   - 一覧取得中に投稿が追加・削除されてもズレやスキップが発生しにくい
// 点。
//
// ドメイン値オブジェクトとして扱うため、フィールドは非公開とし参照はアクセサ経由。
// JSON 永続化のためのデータ変換は postCursorDto を中継して行う。
type PostCursor struct {
	lastPostID     PostID
	lastCreateTime time.Time
}

// NewPostCursor は PostCursor の唯一のコンストラクタ。
func NewPostCursor(lastPostID PostID, lastCreateTime time.Time) *PostCursor {
	return &PostCursor{
		lastPostID:     lastPostID,
		lastCreateTime: lastCreateTime,
	}
}

func (c *PostCursor) LastPostID() PostID {
	return c.lastPostID
}

func (c *PostCursor) LastCreateTime() time.Time {
	return c.lastCreateTime
}

// postCursorDto は PostCursor を JSON で表現するための中間データ。
//
// encoding/json は非公開フィールドを扱えないため、ドメイン値として
// フィールドを隠したい PostCursor と、JSON タグ付きの公開フィールドが
// 必要なシリアライザの橋渡しを担う。
// 本型は純粋なデータ転送オブジェクトで、不変条件や振る舞いは持たせない。
// domain パッケージ内部の都合なので非公開で定義する。
type postCursorDto struct {
	LastPostID     string    `json:"last_post_id"`
	LastCreateTime time.Time `json:"last_create_time"`
}

// Marshal は PostCursor を Pagination.cursor として永続化するための
// バイト列に変換する。
//
// 表現形式に JSON を選んでいる理由は、
//   - 人間が DB を覗いたときにデバッグしやすい
//   - スキーマ変更（フィールド追加）に対して前方互換を取りやすい
//   - protobuf binary などに比べて依存関係が軽い
// のバランス。性能要件が厳しくないコンテキストでは妥当なトレードオフ。
func (c *PostCursor) Marshal() ([]byte, error) {
	return json.Marshal(postCursorDto{
		LastPostID:     c.lastPostID.String(),
		LastCreateTime: c.lastCreateTime,
	})
}

// UnmarshalPostCursor は Marshal の逆変換。
// Pagination から取り出した cursor バイト列を PostCursor に復元する。
func UnmarshalPostCursor(data []byte) (*PostCursor, error) {
	var dto postCursorDto
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}
	return NewPostCursor(NewPostID(dto.LastPostID), dto.LastCreateTime), nil
}
