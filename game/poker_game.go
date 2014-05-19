package game

import (
	"github.com/elricL/poker/board"
	"log"
)

func gameStart(board *board.Board) {
	board.Shuffle()
	var i = board.Dealer
	for {
		var card1 = board.Deck.GetPokerCard()
		var card2 = board.Deck.GetPokerCard()
		sendPokerMessage(card1, i.Conn)
		sendPokerMessage(card2, i.Conn)
		i.Hand = []string{card1, card2}
		sendPokerMessage("ciruclar Nub"+i.Name, i.Conn)
		i = i.Next_player
		if i == board.Dealer {
			break
		}
	}
	log.Println(board.Deck)
	board.GameState = "preFlop"
}
