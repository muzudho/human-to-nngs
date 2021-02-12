package controller

// NngsListener - NNGS からの受信メッセージをさばきます。
type NngsListener interface {
	// MyPhase - 自分の手番であることのアナウンスが来ました。
	// この通知を受け取ったら、思考を開始してください。
	// 指し手の入力をするには、別途、非同期の出力で 返してください
	MyPhase()

	// OpponentPhase - 相手の手番であることのアナウンス
	OpponentPhase()
}
