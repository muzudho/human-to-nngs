package controller

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	e "github.com/muzudho/human-to-nngs/entities"
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

// Spawn - クライアント接続
func Spawn(connectorConf e.ConnectorConf) error {
	// NNGSクライアントの状態遷移図
	nngsClientStateDiagram := NngsClientStateDiagram{
		connectorConf: connectorConf,
		// nngsClientStateDiagram: *new(NngsClientStateDiagram),
		index:                  0,
		regexCommand:           *regexp.MustCompile("^(\\d+) (.*)"),
		regexUseMatch:          *regexp.MustCompile("^Use <match"),
		regexUseMatchToRespond: *regexp.MustCompile("^Use <(.+?)> or <(.+?)> to respond."), // 頭の '9 ' は先に削ってあるから ここに含めない（＾～＾）
		regexMatchAccepted:     *regexp.MustCompile("^Match \\[.+?\\] with (\\S+?) in \\S+? accepted."),
		regexDecline1:          *regexp.MustCompile("declines your request for a match."),
		regexDecline2:          *regexp.MustCompile("You decline the match offer from"),
		regexOneSeven:          *regexp.MustCompile("1 7"),
		regexGame:              *regexp.MustCompile("Game (\\d+) ([a-zA-Z]): (\\S+) \\((\\S+) (\\S+) (\\S+)\\) vs (\\S+) \\((\\S+) (\\S+) (\\S+)\\)"),
		regexMove:              *regexp.MustCompile("\\s*(\\d+)\\(([BWbw])\\): ([A-Z]\\d+|Pass)"),
		regexAcceptCommand:     *regexp.MustCompile("match \\S+ \\S+ (\\d+) ")}
	return telnet.DialToAndCall(fmt.Sprintf("%s:%d", connectorConf.Server.Host, connectorConf.Server.Port), nngsClientStateDiagram)
}

// CallTELNET - 決まった形のメソッド。
func (dia NngsClientStateDiagram) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	print("[情報] 受信開始☆")
	lis := nngsClientStateDiagramListener{}

	dia.writer = w
	dia.reader = r

	go dia.read(&lis)

	// scanner - 標準入力を監視します。
	scanner := bufio.NewScanner(os.Stdin)
	// 無限ループ。 一行読み取ります。
	for scanner.Scan() {
		// 書き込みます。最後に改行を付けます。
		oi.LongWrite(dia.writer, scanner.Bytes())
		oi.LongWrite(dia.writer, []byte("\n"))
	}
}

// 送られてくるメッセージを待ち構えるループです。
func (dia *NngsClientStateDiagram) read(lis *nngsClientStateDiagramListener) {
	var buffer [1]byte // これが満たされるまで待つ。1バイト。
	p := buffer[:]

	for {
		n, err := dia.reader.Read(p) // 送られてくる文字がなければ、ここでブロックされます。

		if n > 0 {
			bytes := p[:n]
			dia.lineBuffer[dia.index] = bytes[0]
			dia.index++

			if dia.newlineReadableState < 2 {
				// [受信] 割り込みで 改行がない行も届くので、改行が届くまで待つという処理ができません。
				print(string(bytes)) // 受け取るたびに１文字ずつ表示。
			}

			// 改行を受け取る前にパースしてしまおう☆（＾～＾）早とちりするかも知れないけど☆（＾～＾）
			dia.parse(lis)

			// `Login:` のように 改行が送られてこないケースはあるが、
			// 対局が始まってしまえば、改行は送られてくると考えろだぜ☆（＾～＾）
			if bytes[0] == '\n' {
				dia.index = 0

				if dia.newlineReadableState == 1 {
					print("[行単位入力へ切替(^q^)]")
					dia.newlineReadableState = 2
					break // for文を抜ける
				}
			}
		}

		if nil != err {
			return // 相手が切断したなどの理由でエラーになるので、終了します。
		}
	}

	// 改行が送られてくるものと考えるぜ☆（＾～＾）
	// これで、１行ずつ読み込めるな☆（＾～＾）
	for {
		n, err := dia.reader.Read(p) // 送られてくる文字がなければ、ここでブロックされます。

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
				dia.parse(lis)
				dia.index = 0

			} else {
				dia.lineBuffer[dia.index] = bytes[0]
				dia.index++
			}
		}
	}
}

// 簡易表示モードに切り替えます。
// Original code: NngsClient.rb/NNGSClient/`def login`
func setClientMode(w telnet.Writer) {
	oi.LongWrite(w, []byte("set client true\n"))
}
