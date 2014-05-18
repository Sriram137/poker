package main

import (
	"github.com/gorilla/websocket"
	"log"
)

func sendPokerMessage(msg string, conn *websocket.Conn) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		delete(connections, conn)
		conn.Close()
	}
}

func changeDealer() {
	dealer = (dealer + 1) % currentPlayers.length()
}

func gameStart() {
	changeDealer()
	expectedPlayer = dealer
	for name, con := range currentPlayers.connMap {
		var card1 = getPokerCard()
		var card2 = getPokerCard()
		sendPokerMessage(card1, con)
		sendPokerMessage(card2, con)
		sendPokerMessage("Hello Nub "+name, con)
	}
	log.Println(cardDeck)
	gameState = "preFlop"
}

func nextTurn() {
	expectedPlayer = (expectedPlayer + 1) % currentPlayers.length()
}
