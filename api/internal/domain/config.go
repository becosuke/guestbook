package domain

import (
	"go.uber.org/zap/zapcore"
)

// Config はアプリケーション起動時に解決される設定値の集約。
//
// あえて adapter/infrastructure 側ではなくドメインに置いているのは、
//   - 全層が依存する設定値の「正規の型」をドメインが定義することで、
//     adapter 層をすげ替えても設定の表現がブレないようにする
//   - 設定値そのもの（環境、ホスト、ポート、DSN）はアプリケーションの
//     振る舞いを決定するドメイン関心ごとに近い
// と判断したため。
//
// zapcore.Level だけは外部パッケージに依存しているが、ログレベルの
// 標準的な表現として広く受け入れられているので、ここは実用性を重視した
// 妥協として割り切っている。
//
// 値の読み出しは adapter/infrastructure/config 側が envconfig 経由で行い、
// このドメイン型に詰め直して各層に注入する。
type Config struct {
	Environment Environment
	LogLevel    zapcore.Level
	GrpcHost    string
	GrpcPort    int
	RestHost    string
	RestPort    int
	DatabaseURL string
}
