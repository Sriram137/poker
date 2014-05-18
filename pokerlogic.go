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
}

func sendAll(msg []byte) {
	for conn := range connections {
		sendPokerMessage(string(msg), conn)
	}
}

func handlePokerMessage(msg []byte, conn *websocket.Conn) {
	fmt.Println(gameState)
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
		if gameState == "waiting" {
			board.addPlayer(Player{nil, false, conn, name, nil})
			currentPlayers.addPlayer(name, conn)
			if currentPlayers.length() > 2 {
				gameState = "canStart"
				gameStart()
			}
			log.Println(commandMsg)
		}
	}
	log.Println(currentPlayers)
	log.Println(string(msg))
	fmt.Println(gameState)
	fmt.Println()
	// sendAll(msg)
}
