package controller

// NngsHumanController - NNGS からの受信メッセージをさばきます。
type NngsHumanController struct {
	NngsListener

	// EntryConf - 参加設定
	EntryConf EntryConf
}

// MyPhase - 自分の手番であることのアナウンスが来ました。
// この通知を受け取ったら、思考を開始してください。
// 指し手の入力をするには、別途、非同期の出力で 返してください
// Original code: nngsCUI.rb/announce class/update/`when 'my_turn'`
func (con NngsHumanController) MyPhase() {
	print("****** I am thinking now   ******")
}

// OpponentPhase - 相手の手番であることのアナウンス
// Original code: nngsCUI.rb/announce class/update/`when 'his_turn'`
func (con NngsHumanController) OpponentPhase() {
	print("****** wating for his move ******")
}
