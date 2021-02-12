package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/muzudho/human-to-nngs/controller/clistat"
	"github.com/muzudho/human-to-nngs/entities/phase"
	"github.com/reiver/go-oi"
)

// NngsClientStateDiagram - NNGSクライアントの状態遷移図
type NngsClientStateDiagram struct {
	// 状態遷移の中の小さな区画
	promptState int
}

func (dia *NngsClientStateDiagram) promptDiagram(lib *libraryListener, subCode int) {
	switch subCode {
	// Info
	case 5:
		if dia.promptState == 7 {
			// 対局終了
			lib.matchEnd()
		}
		dia.promptState = 5
	// PlayingGo
	case 6:
		if dia.promptState == 5 {
			// 対局成立
			lib.matchStart()
		}
		dia.promptState = 6
	// Scoring
	case 7:
		if dia.promptState == 6 {
			// 得点計算
			lib.scoring()
		}
		dia.promptState = 7
	default:
		// "1 1" とか来ても無視しろだぜ☆（＾～＾）
	}
}

func (lib *libraryListener) parse() {
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

			oi.LongWrite(lib.writer, []byte(user))
			oi.LongWrite(lib.writer, []byte("\n"))

			lib.state = clistat.EnteredMyName
		}
	// Original code: NngsClient.rb/NNGSClient/`def login`
	case clistat.EnteredMyName:
		if line == "1 1" {
			// パスワードを入れろだぜ☆（＾～＾）
			if lib.entryConf.Pass() == "" {
				panic("Need password")
			}
			oi.LongWrite(lib.writer, []byte(lib.entryConf.Nngs.Pass))
			oi.LongWrite(lib.writer, []byte("\n"))
			setClientMode(lib.writer)
			lib.state = clistat.EnteredClientMode

		} else if line == "Password: " {
			// パスワードを入れろだぜ☆（＾～＾）
			if lib.entryConf.Pass() == "" {
				panic("Need password")
			}
			oi.LongWrite(lib.writer, []byte(lib.entryConf.Nngs.Pass))
			oi.LongWrite(lib.writer, []byte("\n"))
			lib.state = clistat.EnteredMyPasswordAndIAmWaitingToBePrompted

		} else if line == "#> " {
			setClientMode(lib.writer)
			lib.state = clistat.EnteredClientMode
		}
		// 入力した名前が被っていれば、ここで無限ループしてるかも☆（＾～＾）

	// Original code: NngsClient.rb/NNGSClient/`def login`
	case clistat.EnteredMyPasswordAndIAmWaitingToBePrompted:
		if line == "#> " {
			setClientMode(lib.writer)
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
				fmt.Printf("[情報] 白番として、対局を申し込んだぜ☆（＾～＾）[%s]", message)
				oi.LongWrite(lib.writer, []byte(message))
			case "B", "b":
				lib.MyColor = phase.Black
				message := fmt.Sprintf("match %s B %d %d %d\n", lib.entryConf.Opponent(), lib.entryConf.BoardSize(), lib.entryConf.AvailableTimeMinutes(), lib.entryConf.CanadianTiming())
				fmt.Printf("[情報] 黒番として、対局を申し込んだぜ☆（＾～＾）[%s]", message)
				oi.LongWrite(lib.writer, []byte(message))
			default:
				panic(fmt.Sprintf("Unexpected phase [%s].", lib.entryConf.Phase()))
			}
		}
		lib.state = clistat.WaitingInInfo
	case clistat.WaitingInInfo:
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
			promptStateBytes := matches[2]

			code, err := strconv.Atoi(commandCode)
			if err != nil {
				// 想定外の遷移だぜ☆（＾～＾）！
				panic(err)
			}
			switch code {
			// Prompt
			case 1:
				promptState := string(promptStateBytes)
				promptStateNum, err := strconv.Atoi(promptState)
				if err == nil {
					lib.nngsClientStateDiagram.promptDiagram(lib, promptStateNum)
				}
			// Info
			case 9:
				// parse_9
				// print("[9だぜ☆]")
				if lib.regexUseMatch.Match(promptStateBytes) {
					matches2 := lib.regexUseMatchToRespond.FindSubmatch(promptStateBytes)
					if 2 < len(matches2) {
						// 対局を申し込まれた方だけ、ここを通るぜ☆（＾～＾）
						// Original code: cmd_match_ok
						// 3回ぐらい ここを通るような？.
						fmt.Printf("[情報] 対局が付いたぜ☆（＾～＾）accept[%s],decline[%s]\n", matches2[1], matches2[2])

						// Example: `match kifuwarabi W 19 40 0`
						lib.CommandOfMatchAccept = string(matches2[1])
						// Example: `decline kifuwarabi`
						lib.CommandOfMatchDecline = string(matches2[2])

						// match_request
						// request
						// ask_match
						// puts 'match requested. accept? (Y/n):'
						// コンピューター・プレイヤーなら常に承諾します。
						// message := fmt.Sprintf("%s\n", lib.CommandOfMatchAccept)
						// oi.LongWrite(lib.writer, []byte(message))

						// acceptコマンドを半角空白でスプリットした３番目が、申し込んできた方の手番
						matchAcceptTokens := strings.Split(lib.CommandOfMatchAccept, " ")
						if len(matchAcceptTokens) < 6 {
							panic(fmt.Sprintf("Error matchAcceptTokens=[%s].", matchAcceptTokens))
						}

						opponentPlayerName := matchAcceptTokens[1]

						opponentColor := matchAcceptTokens[2]
						opponentColorUppercase := strings.ToUpper(opponentColor)
						switch opponentColor {
						case "W":
							lib.MyColor = phase.Black
						case "B":
							lib.MyColor = phase.White
						default:
							panic(fmt.Sprintf("Unexpected opponentColor=%s.", opponentColor))
						}
						boardSize, err := strconv.ParseUint(matchAcceptTokens[3], 10, 0)
						if err != nil {
							panic(err)
						}
						lib.BoardSize = uint(boardSize)
						fmt.Printf("[情報] ボードサイズは%d☆（＾～＾）", lib.BoardSize)

						configuredColor := phase.PhaseNone
						switch lib.entryConf.Phase() {
						case "W", "w":
							// Original code: @color = WHITE
							configuredColor = phase.White
						case "B", "b":
							configuredColor = phase.Black
						default:
							panic(fmt.Sprintf("Unexpected phase [%s].", lib.entryConf.Phase()))
						}

						if lib.MyColor != configuredColor {
							panic(fmt.Sprintf("Unexpected phase. lib.MyColor=%d configuredColor=%d.", lib.MyColor, configuredColor))
						}

						// cmd_match
						message := fmt.Sprintf("match %s %s %d %d %d\n", opponentPlayerName, opponentColorUppercase, lib.entryConf.BoardSize(), lib.entryConf.AvailableTimeMinutes(), lib.entryConf.CanadianTiming())
						fmt.Printf("[情報] 対局を申し込むぜ☆（＾～＾）[%s]\n", message)
						oi.LongWrite(lib.writer, []byte(message))
					}
				} else if lib.regexMatchAccepted.Match(promptStateBytes) {
					// 黒の手番から始まるぜ☆（＾～＾）
					lib.Phase = phase.Black

				} else if lib.regexDecline1.Match(promptStateBytes) {
					print("[対局はキャンセルされたぜ☆]")
					// self.match_cancel
				} else if lib.regexDecline2.Match(promptStateBytes) {
					print("[対局はキャンセルされたぜ☆]")
					// self.match_cancel
				} else if lib.regexOneSeven.Match(promptStateBytes) {
					print("[サブ遷移へ☆]")
					lib.nngsClientStateDiagram.promptDiagram(lib, 7)
				} else {
					// "9 1 5" とか来るが、無視しろだぜ☆（＾～＾）
				}
			// Move
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
					matches2 := lib.regexGame.FindSubmatch(promptStateBytes)
					if 10 < len(matches2) {
						// 白 VS 黒 の順序固定なのか☆（＾～＾）？ それともマッチを申し込んだ方 VS 申し込まれた方 なのか☆（＾～＾）？
						fmt.Printf("[情報] 対局現在情報☆（＾～＾） gameid[%s], gametype[%s] white_user[%s][%s][%s][%s] black_user[%s][%s][%s][%s]", matches2[1], matches2[2], matches2[3], matches2[4], matches2[5], matches2[6], matches2[7], matches2[8], matches2[9], matches2[10])

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
					matches2 := lib.regexMove.FindSubmatch(promptStateBytes)
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

						fmt.Printf("[情報] 初回指し手 lib.MyColor=%d, lib.Phase=%d", lib.MyColor, lib.Phase)
						if lib.MyColor == lib.Phase {
							// 自分の手番だぜ☆（＾～＾）！
							lib.OpponentMove = string(matches2[3]) // 相手の指し手が付いてくるので記憶
							fmt.Printf("[情報] ここを通ってるのを見たことはないが、自分の手番で一旦ブロッキング☆（＾～＾）")
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
