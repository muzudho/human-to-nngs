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
	// WaitingInTheLobby - 対局が申し込まれるのをロビーで待ちます
	WaitingInTheLobby
	// BlockingMyTurn - 自分の手番で受信はブロック中です。
	BlockingMyTurn
	// BlockingOpponentTurn - 相手の手番で受信はブロック中です。
	BlockingOpponentTurn
)
