package entities

import (
	"fmt"
	"strings"

	"github.com/muzudho/human-to-nngs/entities/phase"
)

// EntryConf - 参加設定。
type EntryConf struct {
	Server           Server
	User             User
	MatchApplication MatchApplication
}

// Server - [Server] 区画。
type Server struct {
	Host string
	Port int64 // Tomlのライブラリーが精度を細かく指定できないので int64 型で。
}

// User - [User] 区画。
type User struct {
	// InterfaceType - プレイヤーの種類
	// * `Human` - 人間プレイヤーが接続する
	// * `GTP` - GTP(碁テキスト プロトコル)を用いる思考エンジンが接続する
	InterfaceType string
	// User name
	Name                string
	Pass                string
	EngineCommand       string
	EngineCommandOption string
}

// MatchApplication - [MatchApplication] 区画。
type MatchApplication struct {
	ApplyFromMe          bool
	OpponentName         string
	Phase                string
	BoardSize            int64
	AvailableTimeMinutes int64
	CanadianTiming       int64
}

// InterfaceType - プレイヤーの種類
// * `Human` - 人間プレイヤーが接続する
// * `GTP` - GTP(碁テキスト プロトコル)を用いる思考エンジンが接続する
func (config EntryConf) InterfaceType() string {
	return config.User.InterfaceType
}

// Host - 接続先ホスト名
func (config EntryConf) Host() string {
	return config.Server.Host
}

// Port - 接続先ホストのポート番号
func (config EntryConf) Port() uint {
	return uint(config.Server.Port)
}

// UserName - 対局者名（アカウント名）
// Only A-Z a-z 0-9
// Names may be at most 10 characters long
func (config EntryConf) UserName() string {
	return config.User.Name
}

// Pass - 何路盤
func (config EntryConf) Pass() string {
	return config.User.Pass
}

// EngineCommand - 思考エンジンを起動するコマンドの実行ファイル名の部分（OSにより書き方が異なるかも）
func (config EntryConf) EngineCommand() string {
	return config.User.EngineCommand
}

// EngineCommandOption - 思考エンジンを起動するコマンドの半角スペース区切りの引数（OSにより書き方が異なるかも）
func (config EntryConf) EngineCommandOption() string {
	return config.User.EngineCommandOption
}

// ApplyFromMe - 自分の方から申し込むなら true, 申し込みを受けるのを待つ方なら false。
// true にしたなら、 OpponentName も設定してください
func (config EntryConf) ApplyFromMe() bool {
	return config.MatchApplication.ApplyFromMe
}

// OpponentName - 自分の方から申し込むなら、対戦相手のアカウント名も指定してください。そうでないなら無視されます
func (config EntryConf) OpponentName() string {
	return config.MatchApplication.OpponentName
}

// Phase - 自分の色
func (config EntryConf) Phase() string {
	return config.MatchApplication.Phase
}

// MyColor - 自分の石の色
func (config EntryConf) MyColor() (phase.Phase, string) {
	configuredColorUpperCase := strings.ToUpper(config.MatchApplication.Phase)
	myPhase := phase.PhaseNone
	switch configuredColorUpperCase {
	case "W":
		myPhase = phase.White
	case "B":
		myPhase = phase.Black
	default:
		panic(fmt.Sprintf("Unexpected MatchApplication.Phase [%s].", config.MatchApplication.Phase))
	}

	return myPhase, configuredColorUpperCase
}

// BoardSize - 何路盤
func (config EntryConf) BoardSize() uint {
	return uint(config.MatchApplication.BoardSize)
}

// AvailableTimeMinutes - 持ち時間（分）
func (config EntryConf) AvailableTimeMinutes() uint {
	return uint(config.MatchApplication.AvailableTimeMinutes)
}

// CanadianTiming - カナダ式秒読み。25手を何分以内に打てばよいか
func (config EntryConf) CanadianTiming() uint {
	return uint(config.MatchApplication.CanadianTiming)
}
