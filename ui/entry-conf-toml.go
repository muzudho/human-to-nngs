package ui

import (
	"fmt"
	"io/ioutil"

	e "github.com/muzudho/human-to-nngs/entities"
	u "github.com/muzudho/human-to-nngs/usecases"
	"github.com/pelletier/go-toml"
)

// e "github.com/muzudho/kifuwarabe-uec12/entities"

// LoadEntryConf - Toml形式の参加設定ファイルを読み込みます。
func LoadEntryConf(path string) e.EntryConf {

	// ファイル読込
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		u.G.Chat.Fatal("path=%s", path)
		panic(err)
	}

	debugPrintToml(fileData)

	// Toml解析
	binary := []byte(string(fileData))
	config := e.EntryConf{}
	toml.Unmarshal(binary, &config)

	debugPrintConfig(config)

	return config
}

func debugPrintToml(fileData []byte) {
	fmt.Printf("[情報] %s", string(fileData))
	tomlTree, err := toml.Load(string(fileData))
	if err != nil {
		panic(err)
	}
	fmt.Println("[情報] Input:")
	fmt.Printf("[情報] Server.Host=%s\n", tomlTree.Get("Server.Host").(string))
	fmt.Printf("[情報] Server.Port=%d\n", tomlTree.Get("Server.Port").(int64))
	fmt.Printf("[情報] User.Name=%s\n", tomlTree.Get("User.Name").(string))
	fmt.Printf("[情報] User.Pass=%s\n", tomlTree.Get("User.Pass").(string))
	fmt.Printf("[情報] User.InterfaceType=%s\n", tomlTree.Get("User.InterfaceType").(string))
	fmt.Printf("[情報] User.EngineCommand=%s\n", tomlTree.Get("User.EngineCommand").(string))
	fmt.Printf("[情報] User.EngineCommandOption=%s\n", tomlTree.Get("User.EngineCommandOption").(string))
	fmt.Printf("[情報] MatchApplication.Phase=%s\n", tomlTree.Get("MatchApplication.Phase").(string))
	fmt.Printf("[情報] MatchApplication.BoardSize=%d\n", tomlTree.Get("MatchApplication.BoardSize").(int64))
	fmt.Printf("[情報] MatchApplication.AvailableTimeMinutes=%d\n", tomlTree.Get("MatchApplication.AvailableTimeMinutes").(int64))
	fmt.Printf("[情報] MatchApplication.CanadianTiming=%d\n", tomlTree.Get("MatchApplication.CanadianTiming").(int64))
}

func debugPrintConfig(config e.EntryConf) {
	fmt.Println("[情報] Memory:")
	fmt.Printf("[情報] Server.Host=%s\n", config.Server.Host)
	fmt.Printf("[情報] Server.Port=%d\n", config.Server.Port)
	fmt.Printf("[情報] User.Name=%s\n", config.User.Name)
	fmt.Printf("[情報] User.Pass=%s\n", config.User.Pass)
	fmt.Printf("[情報] User.InterfaceType=%s\n", config.User.InterfaceType)
	fmt.Printf("[情報] User.EngineCommand=%s\n", config.User.EngineCommand)
	fmt.Printf("[情報] User.EngineCommandOption=%s\n", config.User.EngineCommandOption)
	fmt.Printf("[情報] MatchApplication.Phase=%s\n", config.MatchApplication.Phase)
	fmt.Printf("[情報] MatchApplication.BoardSize=%d\n", config.MatchApplication.BoardSize)
	fmt.Printf("[情報] MatchApplication.AvailableTimeMinutes=%d\n", config.MatchApplication.AvailableTimeMinutes)
	fmt.Printf("[情報] MatchApplication.CanadianTiming=%d\n", config.MatchApplication.CanadianTiming)
}
