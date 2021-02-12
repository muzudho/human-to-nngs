package controller

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/muzudho/human-to-nngs/controller/clistat"
	"github.com/muzudho/human-to-nngs/entities/phase"
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

// NngsClient - クライアント
type NngsClient struct {
}

// `github.com/reiver/go-telnet` ライブラリーの動作をリスニングします
type libraryListener struct {
	entryConf EntryConf

	// 末尾に改行が付いていると想定していいフェーズ。逆に、そうでない例は `Login:` とか
	newlineReadableState uint

	// NNGSへ書込み
	writer telnet.Writer
	// NNGSへ読込み
	reader telnet.Reader

	// NNGSクライアントの状態遷移図
	nngsClientStateDiagram NngsClientStateDiagram

	// NNGSの動作をリスニングします
	nngsListener NngsListener
	// １行で 1024 byte は飛んでこないことをサーバーと決めておけだぜ☆（＾～＾）
	lineBuffer [1024]byte
	index      uint

	// 状態遷移
	state clistat.Clistat

	// 正規表現
	regexCommand           regexp.Regexp
	regexUseMatch          regexp.Regexp
	regexUseMatchToRespond regexp.Regexp
	regexMatchAccepted     regexp.Regexp
	regexDecline1          regexp.Regexp
	regexDecline2          regexp.Regexp
	regexOneSeven          regexp.Regexp
	regexGame              regexp.Regexp

	// Example: `15 Game 2 I: kifuwarabe (0 2289 -1) vs kifuwarabi (0 2298 -1)`.
	regexMove          regexp.Regexp
	regexAcceptCommand regexp.Regexp

	// MyColor - 自分の手番の色
	MyColor phase.Phase

	// BoardSize - 何路盤。マッチを受け取ったときに確定
	BoardSize uint
	// Phase - これから指す方。局面の手番とは逆になる
	Phase phase.Phase
	// MyMove - 自分の指し手
	MyMove string
	// OpponentMove - 相手の指し手
	OpponentMove string
	// CommandOfMatchAccept - 申し込まれた対局を受け入れるコマンド。人間プレイヤーの入力補助用
	CommandOfMatchAccept string
	// CommandOfMatchDecline - 申し込まれた対局をお断りするコマンド。人間プレイヤーの入力補助用
	CommandOfMatchDecline string
	// GameID - 対局番号☆（＾～＾） 1 から始まる数☆（＾～＾）
	GameID uint
	// GameType - なんだか分からないが少なくとも "I" とか入ってるぜ☆（＾～＾）
	GameType string
	// GameWName - 白手番の対局者アカウント名
	GameWName string
	// GameWField2 - 白手番の２番目のフィールド（用途不明）
	GameWField2 string
	// GameWAvailableSeconds - 白手番の残り時間（秒）
	GameWAvailableSeconds int
	// GameWField4 - 白手番の４番目のフィールド（用途不明）
	GameWField4 string
	// GameBName - 黒手番の対局者アカウント名
	GameBName string
	// GameBField2 - 黒手番の２番目のフィールド（用途不明）
	GameBField2 string
	// GameBAvailableSeconds - 白手番の残り時間（秒）
	GameBAvailableSeconds int
	// GameBField4 - 黒手番の４番目のフィールド（用途不明）
	GameBField4 string
}

// Spawn - クライアント接続
func (client NngsClient) Spawn(entryConf EntryConf, nngsListener NngsListener) error {
	listener := libraryListener{
		entryConf:              entryConf,
		nngsClientStateDiagram: *new(NngsClientStateDiagram),
		nngsListener:           nngsListener,
		index:                  0,
		regexCommand:           *regexp.MustCompile("^(\\d+) (.*)"),
		regexUseMatch:          *regexp.MustCompile("^Use <match"),                         // (2021-02-12) 先頭付近に '9 ' 追加。
		regexUseMatchToRespond: *regexp.MustCompile("^Use <(.+?)> or <(.+?)> to respond."), // (2021-02-11)末尾に $ 追加☆（＾～＾） // (2021-02-12) 先頭付近に '9 ' 追加。
		regexMatchAccepted:     *regexp.MustCompile("^Match \\[.+?\\] with (\\S+?) in \\S+? accepted."),
		regexDecline1:          *regexp.MustCompile("declines your request for a match."),
		regexDecline2:          *regexp.MustCompile("You decline the match offer from"),
		regexOneSeven:          *regexp.MustCompile("1 7"),
		regexGame:              *regexp.MustCompile("Game (\\d+) ([a-zA-Z]): (\\S+) \\((\\S+) (\\S+) (\\S+)\\) vs (\\S+) \\((\\S+) (\\S+) (\\S+)\\)"),
		regexMove:              *regexp.MustCompile("\\s*(\\d+)\\(([BWbw])\\): ([A-Z]\\d+|Pass)"),
		regexAcceptCommand:     *regexp.MustCompile("match \\S+ \\S+ (\\d+) ")}
	return telnet.DialToAndCall(fmt.Sprintf("%s:%d", entryConf.Nngs.Host, entryConf.Nngs.Port), listener)
}

// CallTELNET - 決まった形のメソッド。
func (lib libraryListener) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	print("[情報] 受信開始☆")

	lib.writer = w
	lib.reader = r

	go lib.read()

	// scanner - 標準入力を監視します。
	scanner := bufio.NewScanner(os.Stdin)
	// 無限ループ。 一行読み取ります。
	for scanner.Scan() {
		// 書き込みます。最後に改行を付けます。
		oi.LongWrite(lib.writer, scanner.Bytes())
		oi.LongWrite(lib.writer, []byte("\n"))
	}
}

// 送られてくるメッセージを待ち構えるループです。
func (lib *libraryListener) read() {
	var buffer [1]byte // これが満たされるまで待つ。1バイト。
	p := buffer[:]

	for {
		n, err := lib.reader.Read(p) // 送られてくる文字がなければ、ここでブロックされます。

		if n > 0 {
			bytes := p[:n]
			lib.lineBuffer[lib.index] = bytes[0]
			lib.index++

			if lib.newlineReadableState < 2 {
				// [受信] 割り込みで 改行がない行も届くので、改行が届くまで待つという処理ができません。
				print(string(bytes)) // 受け取るたびに１文字ずつ表示。
			}

			// 改行を受け取る前にパースしてしまおう☆（＾～＾）早とちりするかも知れないけど☆（＾～＾）
			lib.parse()

			// `Login:` のように 改行が送られてこないケースはあるが、
			// 対局が始まってしまえば、改行は送られてくると考えろだぜ☆（＾～＾）
			if bytes[0] == '\n' {
				lib.index = 0

				if lib.newlineReadableState == 1 {
					lib.newlineReadableState = 2
					// break // for文を抜ける
				}
			}
		}

		if nil != err {
			return // 相手が切断したなどの理由でエラーになるので、終了します。
		}
	}

	/*
		// 改行が送られてくるものと考えるぜ☆（＾～＾）
		// これで、１行ずつ読み込めるな☆（＾～＾）
		for {
			n, err := lib.reader.Read(p) // 送られてくる文字がなければ、ここでブロックされます。

			if nil != err {
				return // 相手が切断したなどの理由でエラーになるので、終了します。
			}

			if n > 0 {
				bytes := p[:n]

				if bytes[0] == '\r' {
					// Windows では、 \r\n と続いてくるものと想定します。
					// Linux なら \r はこないものと想定します。
					continue

				} else if bytes[0] == '\n' {
					// `Login:` のように 改行が送られてこないケースはあるが、
					// 対局が始まってしまえば、改行は送られてくると考えろだぜ☆（＾～＾）
					// 1行をパースします
					lib.parse()
					lib.index = 0

				} else {
					lib.lineBuffer[lib.index] = bytes[0]
					lib.index++
				}
			}
		}
		// */
}

// 簡易表示モードに切り替えます。
// Original code: NngsClient.rb/NNGSClient/`def login`
func setClientMode(w telnet.Writer) {
	oi.LongWrite(w, []byte("set client true\n"))
}

func (lib *libraryListener) matchStart() {
	print("[情報] 対局成立だぜ☆")
}
func (lib *libraryListener) matchEnd() {
	print("[情報] 対局終了だぜ☆")
}
func (lib *libraryListener) scoring() {
	print("[情報] 得点計算だぜ☆")
}
