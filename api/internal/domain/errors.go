package domain

import (
	"errors"
)

// ドメイン共通のエラーを sentinel error として定義する。
//
// 設計方針:
//   - usecase / repository 層は、これらをそのまま返すか fmt.Errorf("...: %w", ErrXxx)
//     で文脈情報をラップして返す。
//   - presentation 層では errors.Is(err, domain.ErrXxx) で振り分けて
//     gRPC の status code（NotFound → codes.NotFound 等）にマッピングする。
//   - これによりドメイン層は gRPC や HTTP のステータスコードを一切知らずに
//     済み、トランスポート差し替え時にもエラー表現がブレない。
//
// AIP-193 のエラーモデルに沿って、汎用的かつ意味の取りやすい名前を採用している。
var (
	// ErrInvalidArgument は引数のセマンティクスが不正な場合のエラー。
	// 形式チェックは presentation 層の protovalidate で弾く前提で、
	// ここではドメインルール上の不正（例: 既に削除済みのリソースを更新しようとした）を表す。
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrAlreadyExists はリソースの一意性制約に違反した場合に返す。
	ErrAlreadyExists = errors.New("already exists")
	// ErrNotFound は指定されたリソースが見つからない場合に返す。
	ErrNotFound = errors.New("not found")
	// ErrInvalidData は DB から取り出したデータがドメインの不変条件を
	// 満たさないなど、外部 IO 越しに壊れた値を観測した場合に使う。
	ErrInvalidData = errors.New("returns invalid data")
	// ErrFailedPrecondition は状態遷移の前提条件が満たされていない場合のエラー。
	// 例えば「削除済みリソースに対する更新」など、引数自体は妥当だが
	// 現在のサーバ状態がオペレーションを許さないケース。
	ErrFailedPrecondition = errors.New("failed precondition")
)
