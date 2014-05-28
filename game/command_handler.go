package game

import (
	"encoding/json"
	"fmt"
	"github.com/elricL/poker/board"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type CommandMsg struct {
	command string
	name    string
	value   int
}

func sendAll(board board.Board, msg []byte) {
	var start = board.Dealer
	for ; ; start = start.Next_player {
		sendPokerMessage(string(msg), start.Conn)
	}
}

func sendPokerMessage(msg string, conn *websocket.Conn) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Println("EERROEROERO WHile sending message")
	}
}

func HandlePokerMessage(msg []byte, pokerBoard *board.Board, conn *websocket.Conn) {
	var commandMsg map[string]string
	var err = json.Unmarshal(msg, &commandMsg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(commandMsg)
	var name = commandMsg["name"]
	switch commandMsg["command"] {
	case "debug":
		log.Println(pokerBoard)
		pokerBoard.Print()
	case "join":
		if pokerBoard.GameState == "waiting" {
			pokerBoard.AddPlayer(board.Player{nil, false, conn, name, nil, 0, 500})
			if pokerBoard.Length() > 2 {
				pokerBoard.GameState = "canStart"
				gameStart(pokerBoard)
				log.Println("#########")
				log.Println(pokerBoard)
			}
		}
	case "check":
		log.Println("Check")
		log.Println(pokerBoard)
		log.Println(pokerBoard.CurrentPlayer)
		if pokerBoard.CurrentPlayer.Conn == conn {
			var moneyToCheck = (pokerBoard.CurrentBet - pokerBoard.CurrentPlayer.CurrentBet)
			pokerBoard.CurrentPlayer.Money = pokerBoard.CurrentPlayer.Money - moneyToCheck
			pokerBoard.Pot += moneyToCheck
			pokerBoard.CurrentPlayer.CurrentBet = pokerBoard.CurrentBet
			pokerBoard.CurrentPlayer = pokerBoard.CurrentPlayer.Next_player
			log.Println("CALLS")
			log.Println(moneyToCheck)
			if pokerBoard.CurrentPlayer == pokerBoard.Starter {
				switch pokerBoard.GameState {
				case "preFlop":
					pokerBoard.GameState = "flop"
					goFlopStuff(pokerBoard)
				}
			}
		}
	case "raise":
		log.Println("Raise")
		log.Println(pokerBoard)
		log.Println(pokerBoard.CurrentPlayer)
		if pokerBoard.CurrentPlayer.Conn == conn {
			raiseAmount, _ := strconv.Atoi(commandMsg["value"])
			pokerBoard.CurrentPlayer.CurrentBet = raiseAmount
			difference := (pokerBoard.CurrentBet - pokerBoard.CurrentPlayer.CurrentBet)
			pokerBoard.CurrentPlayer.Money = pokerBoard.CurrentPlayer.Money - difference
			pokerBoard.Starter = pokerBoard.CurrentPlayer
			pokerBoard.CurrentPlayer = pokerBoard.CurrentPlayer.Next_player
			pokerBoard.Pot += difference
		}
	}
	fmt.Println(pokerBoard.GameState)
	fmt.Println()
	// sendAll(msg)
}
