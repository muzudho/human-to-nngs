package controller

// `github.com/reiver/go-telnet` ライブラリーの動作をリスニングします
type nngsClientStateDiagramListener struct {
}

func (lis *nngsClientStateDiagramListener) matchStart() {
	print("[情報] 対局成立だぜ☆")
}
func (lis *nngsClientStateDiagramListener) matchEnd() {
	print("[情報] 対局終了だぜ☆")
}
func (lis *nngsClientStateDiagramListener) scoring() {
	print("[情報] 得点計算だぜ☆")
}

func (lis *nngsClientStateDiagramListener) myTurn(dia *NngsClientStateDiagram) {
	print("****** I am thinking now   ******")
}
func (lis *nngsClientStateDiagramListener) opponentTurn(dia *NngsClientStateDiagram) {
	print("****** wating for his move ******")
}
