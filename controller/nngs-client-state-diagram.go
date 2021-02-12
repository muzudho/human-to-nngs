package controller

// NngsClientStateDiagram - NNGSクライアントの状態遷移図
type NngsClientStateDiagram struct {
}

func (dia *NngsClientStateDiagram) parseSub1(lib *libraryListener, subCode int) {
	switch subCode {
	case 5:
		if lib.stateSub1 == 7 {
			print("[マッチが終わったぜ☆]")
		}
		lib.stateSub1 = 5
	case 6:
		if lib.stateSub1 == 5 {
			print("[手番が変わったぜ☆]")
		}
		lib.stateSub1 = 6
	case 7:
		if lib.stateSub1 == 6 {
			print("[得点計算だぜ☆]")
		}
		lib.stateSub1 = 7
	default:
		// "1 1" とか来ても無視しろだぜ☆（＾～＾）
	}
}
