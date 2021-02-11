package main

import (
	"flag"
	"fmt"

	e "github.com/muzudho/human-to-nngs/entities"
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

}
