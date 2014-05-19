package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type CommandMsg struct {
	command string
	name    string
	value   int
}

func sendAll(msg []byte) {
	for conn := range connections {
		sendPokerMessage(string(msg), conn)
	}
}

func handlePokerMessage(msg []byte, conn *websocket.Conn) {
	fmt.Println(board.gameState)
	var commandMsg map[string]string
	var err = json.Unmarshal(msg, &commandMsg)
	if err == nil {
		fmt.Println(err)
	}
	fmt.Println(commandMsg)
	var name = commandMsg["name"]
	switch commandMsg["command"] {
	case "join":
		log.Println("Some one joined")
		if board.gameState == "waiting" {
			board.addPlayer(Player{nil, false, conn, name, nil})
			if board.length() > 2 {
				board.gameState = "canStart"
				gameStart()
			}
		}
	case "check":
		if board.starter.name == commandMsg["name"] {
			if board.gameState == "preflop" {

			}
		}
	}
	board.print()
	fmt.Println(board.gameState)
	fmt.Println()
	// sendAll(msg)
}
