package controller

import (
	e "github.com/muzudho/human-to-nngs/entities"
)

// NngsHumanController - NNGS からの受信メッセージをさばきます。
type NngsHumanController struct {
	// ConnectorConf - 参加設定
	ConnectorConf e.ConnectorConf
}
