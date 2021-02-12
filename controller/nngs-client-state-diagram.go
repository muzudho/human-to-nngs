package controller

// NngsClientStateDiagram - NNGSクライアントの状態遷移図
type NngsClientStateDiagram struct {
	// 状態遷移の中の小さな区画
	stateSub1 int
}

func (dia *NngsClientStateDiagram) parseSub1(lib *libraryListener, subCode int) {
	switch subCode {
	// Info
	case 5:
		if dia.stateSub1 == 7 {
			// 対局終了
			lib.matchEnd()
		}
		dia.stateSub1 = 5
	// PlayingGo
	case 6:
		if dia.stateSub1 == 5 {
			// 対局成立
			lib.matchStart()
		}
		dia.stateSub1 = 6
	// Scoring
	case 7:
		if dia.stateSub1 == 6 {
			// 得点計算
			lib.scoring()
		}
		dia.stateSub1 = 7
	default:
		// "1 1" とか来ても無視しろだぜ☆（＾～＾）
	}
}
