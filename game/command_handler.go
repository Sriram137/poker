package game

import (
	"encoding/json"
	"fmt"
	"github.com/elricL/poker/board"
	"github.com/gorilla/websocket"
	"log"
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
	case "join":
		if pokerBoard.GameState == "waiting" {
			pokerBoard.AddPlayer(board.Player{nil, false, conn, name, nil, 500, 0})
			if pokerBoard.Length() > 2 {
				pokerBoard.GameState = "canStart"
				gameStart(pokerBoard)
				log.Println(pokerBoard.Dealer.Name)
				log.Println(pokerBoard.Starter.Name)
				log.Println(pokerBoard.CurrentPlayer.Name)
			}
		}
	case "check":
		log.Println(pokerBoard.Dealer.Name)
		log.Println(pokerBoard.Starter.Name)
		log.Println(pokerBoard.CurrentPlayer.Name)
		log.Println(commandMsg["name"])
		if pokerBoard.CurrentPlayer.Name == commandMsg["name"] {
			switch pokerBoard.GameState {
			case "preFlop":
				pokerBoard.CurrentPlayer = pokerBoard.CurrentPlayer.Next_player
				if pokerBoard.CurrentPlayer == pokerBoard.Starter {
					pokerBoard.GameState = "flop"
					goFlopStuff(pokerBoard)
				}
			}
		}
	}
	fmt.Println(pokerBoard.GameState)
	fmt.Println()
	// sendAll(msg)
}
