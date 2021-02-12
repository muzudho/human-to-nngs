package controller

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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
		regexUseMatch:          *regexp.MustCompile("^Use <match"),
		regexUseMatchToRespond: *regexp.MustCompile("^Use <(.+?)> or <(.+?)> to respond.$"), // (2021-02-11)末尾に $ 追加☆（＾～＾）
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

	go lib.read(w, r)

	writeByHuman(w)
}

// 送られてくるメッセージを待ち構えるループです。
func (lib *libraryListener) read(w telnet.Writer, r telnet.Reader) {
	var buffer [1]byte // これが満たされるまで待つ。1バイト。
	p := buffer[:]

	for {
		n, err := r.Read(p) // 送られてくる文字がなければ、ここでブロックされます。

		if n > 0 {
			bytes := p[:n]
			lib.lineBuffer[lib.index] = bytes[0]
			lib.index++
			// [受信] 割り込みで 改行がない行も届くので、改行が届くまで待つという処理ができません。
			print(string(bytes)) // 受け取るたびに１文字ずつ表示。

			// 改行を受け取る前にパースしてしまおう☆（＾～＾）早とちりするかも知れないけど☆（＾～＾）
			lib.parse(w)

			// 行末を判定できるか☆（＾～＾）？
			if bytes[0] == '\n' {
				// print("行末だぜ☆（＾～＾）！")
				// print("<行末☆>")
				lib.index = 0
			}
		}

		if nil != err {
			break // 相手が切断したなどの理由でエラーになるので、終了します。
		}
	}
}

func (lib *libraryListener) parse(w telnet.Writer) {
	// 現在読み取り中の文字なので、早とちりするかも知れないぜ☆（＾～＾）
	line := string(lib.lineBuffer[:lib.index])

	switch lib.state {
	case clistat.None:
		// Original code: NngsClient.rb/NNGSClient/`def login`
		// Waitfor "Login: ".
		if line == "Login: " {
			// あなたの名前を入力してください。

			// 設定ファイルから自動で入力するぜ☆（＾ｑ＾）
			user := lib.entryConf.User()

			// 自動入力のときは、設定ミスなら強制終了しないと無限ループしてしまうぜ☆（＾～＾）
			if user == "" {
				panic("Need name (User)")
			}

			oi.LongWrite(w, []byte(user))
			oi.LongWrite(w, []byte("\n"))

			lib.state = clistat.EnteredMyName
		}
	// Original code: NngsClient.rb/NNGSClient/`def login`
	case clistat.EnteredMyName:
		if line == "1 1" {
			// パスワードを入れろだぜ☆（＾～＾）
			if lib.entryConf.Pass() == "" {
				panic("Need password")
			}
			oi.LongWrite(w, []byte(lib.entryConf.Nngs.Pass))
			oi.LongWrite(w, []byte("\n"))
			setClientMode(w)
			lib.state = clistat.EnteredClientMode

		} else if line == "Password: " {
			// パスワードを入れろだぜ☆（＾～＾）
			if lib.entryConf.Pass() == "" {
				panic("Need password")
			}
			oi.LongWrite(w, []byte(lib.entryConf.Nngs.Pass))
			oi.LongWrite(w, []byte("\n"))
			lib.state = clistat.EnteredMyPasswordAndIAmWaitingToBePrompted

		} else if line == "#> " {
			setClientMode(w)
			lib.state = clistat.EnteredClientMode
		}
		// 入力した名前が被っていれば、ここで無限ループしてるかも☆（＾～＾）

	// Original code: NngsClient.rb/NNGSClient/`def login`
	case clistat.EnteredMyPasswordAndIAmWaitingToBePrompted:
		if line == "#> " {
			setClientMode(w)
			lib.state = clistat.EnteredClientMode
		}
	case clistat.EnteredClientMode:
		if lib.entryConf.Apply() {
			// 対局を申し込みます。
			// 2010/8/25 added by manabe (set color)
			switch lib.entryConf.Phase() {
			case "W", "w":
				// Original code: @color = WHITE
				lib.MyColor = phase.White
				message := fmt.Sprintf("match %s W %d %d %d\n", lib.entryConf.Opponent(), lib.entryConf.BoardSize(), lib.entryConf.AvailableTimeMinutes(), lib.entryConf.CanadianTiming())
				// fmt.Printf("[情報] 対局を申し込んだぜ☆（＾～＾）[%s]", message)
				oi.LongWrite(w, []byte(message))
			case "B", "b":
				lib.MyColor = phase.Black
				message := fmt.Sprintf("match %s B %d %d %d\n", lib.entryConf.Opponent(), lib.entryConf.BoardSize(), lib.entryConf.AvailableTimeMinutes(), lib.entryConf.CanadianTiming())
				// fmt.Printf("[情報] 対局を申し込んだぜ☆（＾～＾）[%s]", message)
				oi.LongWrite(w, []byte(message))
			default:
				panic(fmt.Sprintf("Unexpected phase [%s].", lib.entryConf.Phase()))
			}
		}
		lib.state = clistat.WaitingInTheLobby
	case clistat.WaitingInTheLobby:
		// /^(\d+) (.*)/
		// if lib.regexCommand.MatchString(line) {
		// 	// コマンドの形をしていたぜ☆（＾～＾）
		// 	// fmt.Printf("[情報] 何かコマンドかだぜ☆（＾～＾）？[%s]", line)
		// }
		matches := lib.regexCommand.FindSubmatch(lib.lineBuffer[:lib.index])

		//fmt.Printf("[情報] m[%s]", matches)
		//print(matches)
		if 2 < len(matches) {
			commandCodeBytes := matches[1]
			commandCode := string(commandCodeBytes)
			commandBodyBytes := matches[2]
			// commandBody := string(commandBodyBytes)

			code, err := strconv.Atoi(commandCode)
			if err != nil {
				// 想定外の遷移だぜ☆（＾～＾）！
				panic(err)
			}
			switch code {
			case 1:
				subCode, err := strconv.Atoi(commandCode)
				if err == nil {
					lib.nngsClientStateDiagram.parseSub1(lib, subCode)
				}

			case 9:
				// print("[9だぜ☆]")
				if lib.regexUseMatch.Match(commandBodyBytes) {
					matches2 := lib.regexUseMatchToRespond.FindSubmatch(commandBodyBytes)
					if 2 < len(matches2) {
						// 対局を申し込まれた方だけ、ここを通るぜ☆（＾～＾）
						// Original code: cmd_match_ok
						// 3回ぐらい ここを通るような？.
						fmt.Printf("[情報] 対局が付いたぜ☆（＾～＾）accept[%s],decline[%s]\n", matches2[1], matches2[2])

						// Example: `match kifuwarabi W 19 40 0`
						lib.CommandOfMatchAccept = string(matches2[1])
						// Example: `decline kifuwarabi`
						lib.CommandOfMatchDecline = string(matches2[2])

						// acceptコマンドを半角空白でスプリットした３番目が手番
						myColor := strings.Split(lib.CommandOfMatchAccept, " ")[2]
						switch myColor {
						case "W":
							lib.MyColor = phase.White
						case "B":
							lib.MyColor = phase.Black
						default:
							panic(fmt.Sprintf("Unexpected phase [%s].", myColor))
						}
						// match_accept
						matches3 := lib.regexAcceptCommand.FindSubmatch(commandCodeBytes)
						if 5 < len(matches3) {
							boardSize, err := strconv.ParseUint(string(matches3[1]), 10, 0)
							if err != nil {
								panic(err)
							}
							lib.BoardSize = uint(boardSize)
							fmt.Printf("[情報] ボードサイズは%d☆（＾～＾）", lib.BoardSize)
						}

						switch lib.entryConf.Phase() {
						case "W", "w":
							// Original code: @color = WHITE
							lib.MyColor = phase.White
							message := fmt.Sprintf("match %s W %d %d %d\n", lib.entryConf.Opponent(), lib.entryConf.BoardSize(), lib.entryConf.AvailableTimeMinutes(), lib.entryConf.CanadianTiming())
							fmt.Printf("[情報] 白へ対局を申し込むぜ☆（＾～＾）[%s]\n", message)
							oi.LongWrite(w, []byte(message))
						case "B", "b":
							lib.MyColor = phase.Black
							message := fmt.Sprintf("match %s B %d %d %d\n", lib.entryConf.Opponent(), lib.entryConf.BoardSize(), lib.entryConf.AvailableTimeMinutes(), lib.entryConf.CanadianTiming())
							fmt.Printf("[情報] 黒へ対局を申し込むぜ☆（＾～＾）[%s]\n", message)
							oi.LongWrite(w, []byte(message))
						default:
							panic(fmt.Sprintf("Unexpected phase [%s].", lib.entryConf.Phase()))
						}

						respondToMatchApplication(w, lib.CommandOfMatchAccept, lib.CommandOfMatchDecline)
					}
				} else if lib.regexMatchAccepted.Match(commandBodyBytes) {
					// 黒の手番から始まるぜ☆（＾～＾）
					lib.Phase = phase.Black

				} else if lib.regexDecline1.Match(commandBodyBytes) {
					print("[対局はキャンセルされたぜ☆]")
					// self.match_cancel
				} else if lib.regexDecline2.Match(commandBodyBytes) {
					print("[対局はキャンセルされたぜ☆]")
					// self.match_cancel
				} else if lib.regexOneSeven.Match(commandBodyBytes) {
					print("[サブ遷移へ☆]")
					lib.nngsClientStateDiagram.parseSub1(lib, 7)
				} else {
					// "9 1 5" とか来るが、無視しろだぜ☆（＾～＾）
				}
			// マッチ確立の合図を得たときだぜ☆（＾～＾）
			// Original code: NngsClient.rb/NNGSClient/`def parse_15(code, line)`
			// Example: `15 Game 2 I: kifuwarabe (0 2289 -1) vs kifuwarabi (0 2298 -1)`.
			// Example: `15   4(B): J4`.
			case 15:
				// print("15だぜ☆")
				doing := true

				// 対局中、ゲーム情報は 指し手の前に毎回流れてくるぜ☆（＾～＾）
				// 自分が指すタイミングと、相手が指すタイミングのどちらでも流れてくるぜ☆（＾～＾）
				// とりあえずゲーム情報を全部変数に入れとけばあとで使える☆（＾～＾）
				if doing {
					matches2 := lib.regexGame.FindSubmatch(commandBodyBytes)
					if 10 < len(matches2) {
						// 白 VS 黒 の順序固定なのか☆（＾～＾）？ それともマッチを申し込んだ方 VS 申し込まれた方 なのか☆（＾～＾）？
						// fmt.Printf("[情報] 対局現在情報☆（＾～＾） gameid[%s], gametype[%s] white_user[%s][%s][%s][%s] black_user[%s][%s][%s][%s]", matches2[1], matches2[2], matches2[3], matches2[4], matches2[5], matches2[6], matches2[7], matches2[8], matches2[9], matches2[10])

						// ゲームID
						// Original code: @gameid
						gameID, err := strconv.ParseUint(string(matches2[1]), 10, 0)
						if err != nil {
							panic(err)
						}
						lib.GameID = uint(gameID)

						// ゲームの型？
						// Original code: @gametype
						lib.GameType = string(matches2[2])

						// 白手番の名前、フィールド２、残り時間（秒）、フィールド４
						// Original code: @white_user = [$3, $4, $5, $6]
						lib.GameWName = string(matches2[3])
						lib.GameWField2 = string(matches2[4])

						gameWAvailableSeconds, err := strconv.Atoi(string(matches2[5]))
						if err != nil {
							panic(err)
						}
						lib.GameWAvailableSeconds = gameWAvailableSeconds

						lib.GameWField4 = string(matches2[6])

						// 黒手番の名前、フィールド２、残り時間（秒）、フィールド４
						// Original code: @black_user = [$7, $8, $9, $10]
						lib.GameBName = string(matches2[7])
						lib.GameBField2 = string(matches2[8])

						gameBAvailableSeconds, err := strconv.Atoi(string(matches2[9]))
						if err != nil {
							panic(err)
						}
						lib.GameBAvailableSeconds = gameBAvailableSeconds

						lib.GameBField4 = string(matches2[10])

						doing = false
					}
				}

				// 指し手はこっちだぜ☆（＾～＾）
				if doing {
					matches2 := lib.regexMove.FindSubmatch(commandBodyBytes)
					if 3 < len(matches2) {
						// Original code: @lastmove = [$1, $2, $3]
						fmt.Printf("[情報] 指し手☆（＾～＾） code[%s], color[%s] move[%s]", matches2[1], matches2[2], matches2[3])

						// 相手の指し手を受信したのだから、手番はその逆だぜ☆（＾～＾）
						switch string(matches2[2]) {
						case "B":
							lib.Phase = phase.White
						case "W":
							lib.Phase = phase.Black
						default:
							panic(fmt.Sprintf("Unexpected phase %s", string(matches2[2])))
						}

						if lib.MyColor == lib.Phase {
							// 自分の手番だぜ☆（＾～＾）！
							lib.OpponentMove = string(matches2[3]) // 相手の指し手が付いてくるので記憶
							fmt.Printf("[情報] 自分の手番で一旦ブロッキング☆（＾～＾）")
							// 初回だけここを通るが、以後、ここには戻ってこないぜ☆（＾～＾）
							lib.state = clistat.BlockingMyTurn

							// Original code: nngsCUI.rb/announce class/update/`when 'my_turn'`.
							// Original code: nngsCUI.rb/engine  class/update/`when 'my_turn'`.
							lib.nngsListener.MyPhase()

							// @gtp.time_left('WHITE', @nngs.white_user[2])
							// @gtp.time_left('BLACK', @nngs.black_user[2])
							/*
							   mv, c = @gtp.genmove
							   if mv.nil?
							     mv = 'PASS'
							   elsif mv == "resign"

							   else
							     i, j = mv
							     mv = '' << 'ABCDEFGHJKLMNOPQRST'[i-1]
							     mv = "#{mv}#{j}"
							   end
							   @nngs.input mv
							*/
						} else {
							// 相手の手番だぜ☆（＾～＾）！
							lib.MyMove = string(matches2[3]) // 自分の指し手が付いてくるので記憶
							fmt.Printf("[情報] 相手の手番で一旦ブロッキング☆（＾～＾）")
							// 初回だけここを通るが、以後、ここには戻ってこないぜ☆（＾～＾）
							lib.state = clistat.BlockingOpponentTurn

							// Original code: nngsCUI.rb/annouce class/update/`when 'his_turn'`.
							// Original code: nngsCUI.rb/engine  class/update/`when 'his_turn'`.
							lib.nngsListener.OpponentPhase()

							// lib.
							/*
								      mv = if move == 'Pass'
								             nil
								           elsif move.downcase[/resign/] == "resign"
								             "resign"
								           else
								             i = move.upcase[0].ord - ?A.ord + 1
									         i = i - 1 if i > ?I.ord - ?A.ord
								             j = move[/[0-9]+/].to_i
								             [i, j]
								           end
								#      p [mv, @his_color]
								      @gtp.playmove([mv, @his_color])
							*/
						}

						doing = false
					}
				}
			default:
				// 想定外のコードが来ても無視しろだぜ☆（＾～＾）
			}
		}
	case clistat.BlockingMyTurn:
		// 自分の手番で受信はブロック中です
		// fmt.Printf("[情報] 自分[%d]のターン☆（＾～＾）", lib.MyColor)
	case clistat.BlockingOpponentTurn:
		// 相手の手番で受信はブロック中です。
		// fmt.Printf("[情報] 自分[%d]の相手のターン☆（＾～＾）", lib.MyColor)
	default:
		// 想定外の遷移だぜ☆（＾～＾）！
		panic(fmt.Sprintf("Unexpected state transition. state=%d", lib.state))
	}
}

// 簡易表示モードに切り替えます。
// Original code: NngsClient.rb/NNGSClient/`def login`
func setClientMode(w telnet.Writer) {
	oi.LongWrite(w, []byte("set client true\n"))
}

// 申込みに応えます。
// Original code: match_request(), ask_match().
func respondToMatchApplication(w telnet.Writer, accept string, decline string) {
	// 人間プレイヤーなら、尋ねて応答を待ちます。
	// 'match requested. accept? (Y/n):'
	// if no
	//   match_cancel
	// else
	//   match_ok

	// コンピューター・プレイヤーなら常に承諾します。
	message := fmt.Sprintf("%s\n", accept)
	oi.LongWrite(w, []byte(message))
}

// 人間がいつでも書き込んで送信できるようにするループです。
func writeByHuman(w telnet.Writer) {
	// scanner - 標準入力を監視します。
	scanner := bufio.NewScanner(os.Stdin)
	// 一行読み取ります。
	for scanner.Scan() {
		// 書き込みます。最後に改行を付けます。
		oi.LongWrite(w, scanner.Bytes())
		oi.LongWrite(w, []byte("\n"))
	}
}

func (lib *libraryListener) matchStart() {
	print("[情報] 手番が変わったぜ☆")
}
func (lib *libraryListener) matchEnd() {
	print("[情報] マッチが終わったぜ☆")
}
func (lib *libraryListener) scoring() {
	print("[情報] 得点計算だぜ☆")
}
