package entities

// GlobalVariables - グローバル変数。
type GlobalVariables struct {
	// Log - ロガー。
	Log Logger
	// Chat - チャッター。 標準出力とロガーを一緒にしただけです。
	Chat Chatter
}

// G - グローバル変数。思い切った名前。
var G GlobalVariables
