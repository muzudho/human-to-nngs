package phase

import "fmt"

// Phase - 黒手番、または白手番。
type Phase int

// state
const (
	// None - 開始。
	PhaseNone Phase = iota
	// Black - 自分のアカウント名を入力しました
	Black
	// White - 自分のパスワードを入力し、そしてプロンプトを待っています
	White
)

// FlipColorString - 色を反転
func FlipColorString(color string) string {
	switch color {
	case "B":
		return "W"
	case "W":
		return "B"
	case "b":
		return "w"
	case "w":
		return "b"
	default:
		return color
	}
}

// ToString - 色を大文字アルファベットに変換
func ToString(phase Phase) string {
	switch phase {
	case Black:
		return "B"
	case White:
		return "W"
	default:
		panic(fmt.Sprintf("Unexpected phase=[%d]", phase))
	}
}

// ToNum - アルファベットを色に変換
func ToNum(color string) Phase {
	switch color {
	case "B":
		return Black
	case "W":
		return White
	case "b":
		return Black
	case "w":
		return White
	default:
		panic(fmt.Sprintf("Unexpected color=[%s]", color))
	}
}
