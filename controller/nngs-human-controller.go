package controller

import e "github.com/muzudho/human-to-nngs/entities"

// NngsHumanController - NNGS からの受信メッセージをさばきます。
type NngsHumanController struct {
	e.NngsListener

	// EntryConf - 参加設定
	EntryConf EntryConf
}

// NoticeMyPhase - 自分の手番であることのアナウンスが来ました。
// この通知を受け取ったら、思考を開始してください。
// 指し手の入力をするには、別途、非同期の出力で 返してください
// Original code: nngsCUI.rb/announce class/update/`when 'my_turn'`
func (con NngsHumanController) NoticeMyPhase() {
	print("****** I am thinking now   ******")
}

// NoticeOpponentPhase - 相手の手番であることのアナウンス
// Original code: nngsCUI.rb/announce class/update/`when 'his_turn'`
func (con NngsHumanController) NoticeOpponentPhase() {
	print("****** wating for his move ******")
}
