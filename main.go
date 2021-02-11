package main

import (
	"flag"
	"fmt"

	c "github.com/muzudho/human-to-nngs/controller"
	e "github.com/muzudho/human-to-nngs/entities"
	"github.com/muzudho/human-to-nngs/gateway"
	"github.com/muzudho/human-to-nngs/ui"
)

func main() {
	// コマンドライン引数
	entryConfPath := flag.String("entry", "./input/default.entryConf.toml", "*.entryConf.toml file path.")
	flag.Parse()
	fmt.Println(flag.Args())
	fmt.Printf("entryConfPath=%s", *entryConfPath)

	// グローバル変数の作成
	e.G = *new(e.GlobalVariables)

	// ロガーの作成。
	e.G.Log = *e.NewLogger(
		"output/trace.log",
		"output/debug.log",
		"output/info.log",
		"output/notice.log",
		"output/warn.log",
		"output/error.log",
		"output/fatal.log",
		"output/print.log")

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	e.G.Chat = *e.NewChatter(e.G.Log)

	// 標準出力への表示と、ログへの書き込みを同時に行います。
	// e.G.Chat.Trace("Author: %s\n", e.Author)

	// fmt.Println("設定ファイルを読み込んだろ☆（＾～＾）")
	entryConf := ui.LoadEntryConf(*entryConfPath) // "./input/default.entryConf.toml"

	// NNGSからのメッセージ受信に対応するプログラムを指定したろ☆（＾～＾）
	var nngsController e.NngsListener = nil
	fmt.Printf("(^q^) プレイヤーのタイプ☆ [%s]", entryConf.Nngs.PlayerType)
	// Human と決め打ち
	nngsController = c.NngsHumanController{EntryConf: entryConf}

	fmt.Println("(^q^) 何か文字を打てだぜ☆ 終わりたかったら [Ctrl]+[C]☆")
	nngsClient := gateway.NngsClient{}
	nngsClient.Spawn(entryConf, nngsController)
	fmt.Println("(^q^) おわり☆！")

}
