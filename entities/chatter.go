package entities

import (
	"fmt"
	"net"
)

// Chatter - チャッター。 標準出力とロガーを一緒にしただけです。
type Chatter struct {
	logger Logger
}

// NewChatter - チャッターを作成します。
func NewChatter(logger Logger) *Chatter {
	chatter := new(Chatter)
	chatter.logger = logger
	return chatter
}

// Trace - 本番運用時にはソースコードにも残っていないような内容を書くのに使います。
func (chatter Chatter) Trace(text string, args ...interface{}) {
	fmt.Printf(text, args...)           // 標準出力
	chatter.logger.Trace(text, args...) // ログ
}

// Debug - 本番運用時にもデバッグを取りたいような内容を書くのに使います。
func (chatter Chatter) Debug(text string, args ...interface{}) {
	fmt.Printf(text, args...)           // 標準出力
	chatter.logger.Debug(text, args...) // ログ
}

// Info - 多めの情報を書くのに使います。
func (chatter Chatter) Info(text string, args ...interface{}) {
	fmt.Printf(text, args...)          // 標準出力
	chatter.logger.Info(text, args...) // ログ
}

// Notice - 定期的に動作確認を取りたいような、節目、節目の重要なポイントの情報を書くのに使います。
func (chatter Chatter) Notice(text string, args ...interface{}) {
	fmt.Printf(text, args...)            // 標準出力
	chatter.logger.Notice(text, args...) // ログ
}

// Warn - ハードディスクの残り容量が少ないなど、当面は無視できるが対応はしたいような情報を書くのに使います。
func (chatter Chatter) Warn(text string, args ...interface{}) {
	fmt.Printf(text, args...)          // 標準出力
	chatter.logger.Warn(text, args...) // ログ
}

// Error - 動作不良の内容や、理由を書くのに使います。
func (chatter Chatter) Error(text string, args ...interface{}) {
	fmt.Printf(text, args...)           // 標準出力
	chatter.logger.Error(text, args...) // ログ
}

// Fatal - 強制終了したことを伝えます。
func (chatter Chatter) Fatal(text string, args ...interface{}) {
	fmt.Printf(text, args...)           // 標準出力
	chatter.logger.Fatal(text, args...) // ログ
}

// Print - 必ず出力します。
func (chatter Chatter) Print(text string, args ...interface{}) {
	fmt.Printf(text, args...)           // 標準出力
	chatter.logger.Print(text, args...) // ログ
}

// Send - メッセージを送信します。
func (chatter Chatter) Send(conn net.Conn, text string, args ...interface{}) {
	_, err := fmt.Fprintf(conn, text, args...) // 出力先指定
	if err != nil {
		panic(err)
	}

	fmt.Printf(text, args...)           // 標準出力
	chatter.logger.Print(text, args...) // ログ
}
