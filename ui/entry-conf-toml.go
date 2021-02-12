package ui

import (
	"fmt"
	"io/ioutil"

	c "github.com/muzudho/human-to-nngs/controller"
	u "github.com/muzudho/human-to-nngs/usecases"
	"github.com/pelletier/go-toml"
)

// e "github.com/muzudho/kifuwarabe-uec12/entities"

// LoadEntryConf - Toml形式の参加設定ファイルを読み込みます。
func LoadEntryConf(path string) c.EntryConf {

	// ファイル読込
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		u.G.Chat.Fatal("path=%s", path)
		panic(err)
	}

	debugPrintToml(fileData)

	// Toml解析
	binary := []byte(string(fileData))
	config := c.EntryConf{}
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
	fmt.Printf("[情報] Nngs.Host=%s\n", tomlTree.Get("Nngs.Host").(string))
	fmt.Printf("[情報] Nngs.Port=%d\n", tomlTree.Get("Nngs.Port").(int64))
	fmt.Printf("[情報] Nngs.User=%s\n", tomlTree.Get("Nngs.User").(string))
	fmt.Printf("[情報] Nngs.Pass=%s\n", tomlTree.Get("Nngs.Pass").(string))
	fmt.Printf("[情報] MatchApplication.Phase=%s\n", tomlTree.Get("MatchApplication.Phase").(string))
	fmt.Printf("[情報] MatchApplication.BoardSize=%d\n", tomlTree.Get("MatchApplication.BoardSize").(int64))
	fmt.Printf("[情報] MatchApplication.AvailableTimeMinutes=%d\n", tomlTree.Get("MatchApplication.AvailableTimeMinutes").(int64))
	fmt.Printf("[情報] MatchApplication.CanadianTiming=%d\n", tomlTree.Get("MatchApplication.CanadianTiming").(int64))
}

func debugPrintConfig(config c.EntryConf) {
	fmt.Println("[情報] Memory:")
	fmt.Printf("[情報] Nngs.Host=%s\n", config.Nngs.Host)
	fmt.Printf("[情報] Nngs.Port=%d\n", config.Nngs.Port)
	fmt.Printf("[情報] Nngs.User=%s\n", config.Nngs.User)
	fmt.Printf("[情報] Nngs.Pass=%s\n", config.Nngs.Pass)
	fmt.Printf("[情報] MatchApplication.Phase=%s\n", config.MatchApplication.Phase)
	fmt.Printf("[情報] MatchApplication.BoardSize=%d\n", config.MatchApplication.BoardSize)
	fmt.Printf("[情報] MatchApplication.AvailableTimeMinutes=%d\n", config.MatchApplication.AvailableTimeMinutes)
	fmt.Printf("[情報] MatchApplication.CanadianTiming=%d\n", config.MatchApplication.CanadianTiming)
}
