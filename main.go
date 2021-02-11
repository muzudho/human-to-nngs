package main

import (
	"flag"
	"fmt"
)

func main() {
	// コマンドライン引数
	entryConfPath := flag.String("entry", "./input/default.entryConf.toml", "*.entryConf.toml file path.")
	flag.Parse()
	fmt.Println(flag.Args())
	fmt.Printf("entryConfPath=%s", *entryConfPath)

}
