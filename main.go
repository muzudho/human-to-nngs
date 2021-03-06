package main

import (
	"flag"
	"fmt"

	c "github.com/muzudho/human-to-nngs/controller"
	"github.com/muzudho/human-to-nngs/ui"
	u "github.com/muzudho/human-to-nngs/usecases"
)

func main() {
	// コマンドライン引数
	connectorConfPath := flag.String("connector", "./input/connector.conf.toml", "connector.conf.toml file path.")
	flag.Parse()
	fmt.Printf("[情報] flag.Args()=%s\n", flag.Args())
	fmt.Printf("[情報] connectorConfPath=%s\n", *connectorConfPath)

	// グローバル変数の作成
	u.G = *new(u.GlobalVariables)

	// ロガーの作成。
	u.G.Log = *u.NewLogger(
		"output/trace.log",
		"output/debug.log",
		"output/info.log",
		"output/notice.log",
		"output/warn.log",
		"output/error.log",
		"output/fatal.log",
		"output/print.log")

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	u.G.Chat = *u.NewChatter(u.G.Log)

	// fmt.Println("[情報] 設定ファイルを読み込んだろ☆（＾～＾）")
	connectorConf := ui.LoadConnectorConf(*connectorConfPath) // "./input/connector.conf.toml"

	// NNGSからのメッセージ受信に対応するプログラムを指定したろ☆（＾～＾）
	fmt.Printf("[情報] (^q^) プレイヤーのタイプ☆ [%s]", connectorConf.User.InterfaceType)

	fmt.Println("[情報] (^q^) 何か文字を打てだぜ☆ 終わりたかったら [Ctrl]+[C]☆")
	c.Spawn(connectorConf)
	fmt.Println("[情報] (^q^) おわり☆！")
}
