package ui

import (
	"io/ioutil"

	c "github.com/muzudho/human-to-nngs/controller"
	e "github.com/muzudho/human-to-nngs/entities"
	"github.com/pelletier/go-toml"
)

// e "github.com/muzudho/kifuwarabe-uec12/entities"

// LoadEntryConf - Toml形式の参加設定ファイルを読み込みます。
func LoadEntryConf(path string) c.EntryConf {

	// ファイル読込
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		e.G.Chat.Fatal("path=%s", path)
		panic(err)
	}

	/*
		fmt.Print(string(fileData))
		tomlTree, err := toml.Load(string(fileData))
		if err != nil {
			panic(err)
		}
		fmt.Println("Input:")
		fmt.Printf("Nngs.Host=%s\n", tomlTree.Get("Nngs.Host").(string))
		fmt.Printf("Nngs.Port=%d\n", tomlTree.Get("Nngs.Port").(int64))
		fmt.Printf("Nngs.User=%s\n", tomlTree.Get("Nngs.User").(string))
		fmt.Printf("Nngs.Pass=%s\n", tomlTree.Get("Nngs.Pass").(string))
		fmt.Printf("MatchApplication.Phase=%s\n", tomlTree.Get("MatchApplication.Phase").(string))
		fmt.Printf("MatchApplication.BoardSize=%d\n", tomlTree.Get("MatchApplication.BoardSize").(int64))
		fmt.Printf("MatchApplication.AvailableTimeMinutes=%d\n", tomlTree.Get("MatchApplication.AvailableTimeMinutes").(int64))
		fmt.Printf("MatchApplication.CanadianTiming=%d\n", tomlTree.Get("MatchApplication.CanadianTiming").(int64))
	*/

	// Toml解析
	binary := []byte(string(fileData))
	config := c.EntryConf{}
	toml.Unmarshal(binary, &config)

	/*
		fmt.Println("Memory:")
		fmt.Println("Nngs.Host=", config.Nngs.Host)
		fmt.Println("Nngs.Port=", config.Nngs.Port)
		fmt.Println("Nngs.User=", config.Nngs.User)
		fmt.Println("Nngs.Pass=", config.Nngs.Pass)
		fmt.Println("MatchApplication.Phase=", config.MatchApplication.Phase)
		fmt.Println("MatchApplication.BoardSize=", config.MatchApplication.BoardSize)
		fmt.Println("MatchApplication.AvailableTimeMinutes=", config.MatchApplication.AvailableTimeMinutes)
		fmt.Println("MatchApplication.CanadianTiming=", config.MatchApplication.CanadianTiming)
	*/

	return config
}
