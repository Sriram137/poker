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

func gameStart() {
	board.deck = Deck{}
	board.deck.makeShuffledCardPack()
	var i = board.dealer
	for {
		var card1 = board.deck.getPokerCard()
		var card2 = board.deck.getPokerCard()
		sendPokerMessage(card1, i.conn)
		sendPokerMessage(card2, i.conn)
		i.hand = []string{card1, card2}
		sendPokerMessage("ciruclar Nub"+i.name, i.conn)
		i = i.next_player
		if i == board.dealer {
			break
		}
	}
	log.Println(board.deck)
	board.gameState = "preFlop"
}
