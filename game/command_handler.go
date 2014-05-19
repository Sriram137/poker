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
}

func HandlePokerMessage(msg []byte, pokerBoard *board.Board, conn *websocket.Conn) {
	fmt.Println(pokerBoard.GameState)
	var commandMsg map[string]string
	var err = json.Unmarshal(msg, &commandMsg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(commandMsg)
	var name = commandMsg["name"]
	switch commandMsg["command"] {
	case "join":
		log.Println("Some one joined")
		if pokerBoard.GameState == "waiting" {
			pokerBoard.AddPlayer(board.Player{nil, false, conn, name, nil})
			if pokerBoard.Length() > 2 {
				pokerBoard.GameState = "canStart"
				gameStart(pokerBoard)
			}
		}
	case "check":
		if pokerBoard.Starter.Name == commandMsg["name"] {
			if pokerBoard.GameState == "preflop" {

			}
		}
	}
	pokerBoard.Print()
	fmt.Println(pokerBoard.GameState)
	fmt.Println()
	// sendAll(msg)
}
