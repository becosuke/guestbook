package domain

import (
	"time"

	"github.com/google/uuid"
)

// PostID は Post の識別子を表す値オブジェクト。
//
// 内部表現として UUID を採用しているが、型エイリアスではなく独立した型に
// することで、生の string や uuid.UUID と取り違える事故を防いでいる
// （いわゆる primitive obsession の回避）。
type PostID uuid.UUID

// NewPostID は文字列から PostID を構築する唯一のコンストラクタ。
//
// 主なユースケースはネットワーク境界（gRPC / REST から渡された UUID 文字列）と
// repository 層（DB から取得した PostId カラム）からの再構成。
//
// 不正な UUID 文字列が渡された場合 uuid.MustParse は panic する。
// これは presentation 層の protovalidate インターセプタが事前に形式を弾く前提で、
// ドメインに重複したバリデーションを置かないという方針による割り切り。
//
// 空文字を受け取ったときは uuid.Nil 相当のゼロ値 PostID を返す。これは
// 「ID 未指定」の入力（例えば Create リクエストのリソース名なし）を素直に
// 表現するためで、サーバ側で改めて採番する流れと整合する。
func NewPostID(postID string) PostID {
	if postID == "" {
		return PostID(uuid.Nil)
	}
	return PostID(uuid.MustParse(postID))
}

// String は PostID を UUID 文字列に戻す。
// API レスポンスや DB の WHERE 句に渡すときの正規表現として用いる。
func (p PostID) String() string {
	return uuid.UUID(p).String()
}

// Post はゲストブックの投稿を表すドメインエンティティ。
//
// フィールドはすべて非公開とし、アクセサ経由でのみ読めるようにすることで、
// ドメイン外（usecase / adapter 層）からの不変条件の破壊を防いでいる。
// 生成は NewPost コンストラクタに集約し、Post 自身が「正しい状態」しか
// 取り得ないようにする値オブジェクト指向の設計。
type Post struct {
	postID   PostID
	postBody PostBody
	// previousBody は直前の本文を 1 世代だけ保持する。
	// AIP-148 的な「変更前後の値を返したい」要求と、無制限な履歴テーブルを
	// 設けない方針との折衷案として、更新時に旧 PostBody をこのカラムへ
	// 退避させる戦略を取っている（実際の退避は repository 層の UPDATE 文で実施）。
	previousBody PostBody
	createTime   time.Time
	updateTime   time.Time
	// deleteTime は論理削除のタイムスタンプ。ゼロ値であれば未削除。
	// 物理削除ではなく論理削除を採用しているのは、外部参照（例えば pagination
	// カーソルが指している postID）が突然消えても破綻しないようにするため。
	deleteTime time.Time
}

// NewPost は Post の全フィールドを受け取る基本コンストラクタ。
//
// 主に repository 層が DB から読み出した行を Post に再構成するために使う。
// 引数の検証（空文字、長さ、UUID 形式など）はここでは行わず、
// presentation 層の protovalidate インターセプタおよび個々の値オブジェクト
// （PostID / PostBody）側に責務を寄せている。
// ドメインに重複したバリデーションを持たない方針なので、不正な値を渡したときの
// 振る舞いは「呼び出し側の責任」と割り切っている。
//
// 新規作成シナリオでは引数のほとんどがゼロ値になるため、専用ファクトリ
// CreatePost の利用を推奨する。
func NewPost(postID PostID, postBody PostBody, previousBody PostBody, createTime time.Time, updateTime time.Time, deleteTime time.Time) *Post {
	return &Post{
		postID:       postID,
		postBody:     postBody,
		previousBody: previousBody,
		createTime:   createTime,
		updateTime:   updateTime,
		deleteTime:   deleteTime,
	}
}

// CreatePost は新規作成シナリオ専用のファクトリ。
//
// クライアントが指定するのは本文のみで、それ以外の項目は次のように扱う:
//   - PostID はサーバ側で UUID を採番（NewPostID 経由でコンストラクタ規律に従う）
//   - previousBody は履歴なしのため空
//   - createTime / updateTime / deleteTime は DB 側のタイムスタンプ採番に委ねるためゼロ値
//
// これにより usecase.Create 側で time.Time{} のゼロ値や空 PostBody を並べる
// 必要がなくなり、呼び出し意図（「新規作成」）が一目で分かるようになる。
// 内部では NewPost を呼んでおり、Post のインスタンス化経路は NewPost に
// 一本化する原則を崩していない。
func CreatePost(postBody PostBody) *Post {
	return NewPost(NewPostID(uuid.NewString()), postBody, NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
}

func (p *Post) PostID() PostID {
	return p.postID
}

func (p *Post) PostBody() PostBody {
	return p.postBody
}

func (p *Post) PreviousBody() PostBody {
	return p.previousBody
}

func (p *Post) CreateTime() time.Time {
	return p.createTime
}

func (p *Post) UpdateTime() time.Time {
	return p.updateTime
}

func (p *Post) DeleteTime() time.Time {
	return p.deleteTime
}

// Valid は論理削除されていない（= 生きている）投稿かどうかを返す。
// deleteTime のゼロ値判定をドメイン側にカプセル化することで、
// 「論理削除をどう表現するか」という実装詳細を呼び出し側に漏らさない。
func (p *Post) Valid() bool {
	return p.deleteTime.IsZero()
}
