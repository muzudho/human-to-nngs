package clistat

// Clistat - Client state. NNGSサーバーでゲームをしているクライアントの状態遷移。
type Clistat int

// state
const (
	// None - 開始。
	None Clistat = iota
	// EnteredMyName - 自分のアカウント名を入力しました
	EnteredMyName
	// EnteredMyPasswordAndIAmWaitingToBePrompted - 自分のパスワードを入力し、そしてプロンプトを待っています
	EnteredMyPasswordAndIAmWaitingToBePrompted
	// EnteredClientMode - 簡易表示モードに設定しました
	EnteredClientMode
	// WaitingInInfo - 対局が申し込まれるのを '1 5 Waiting' (ロビー相当)で待ちます
	WaitingInInfo
	// BlockingReceiver - 申し込まれた方はブロック中です。
	BlockingReceiver
	// BlockingSender - 申し込んだ方はブロック中です。
	BlockingSender
)
