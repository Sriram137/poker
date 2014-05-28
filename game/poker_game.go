package game

import (
	"github.com/elricL/poker/board"
	"strconv"
)

func gameStart(pokerBoard *board.Board) {
	pokerBoard.Shuffle()
	var i = pokerBoard.Dealer
	var count = 0
	for {
		var card1 = pokerBoard.Deck.GetPokerCard()
		var card2 = pokerBoard.Deck.GetPokerCard()
		sendPokerMessage(card1, i.Conn)
		sendPokerMessage(card2, i.Conn)
		sendPokerMessage(strconv.Itoa(count), i.Conn)
		i.Hand = []string{card1, card2}
		sendPokerMessage("ciruclar Nub"+i.Name, i.Conn)
		i = i.Next_player
		count++
		if i == pokerBoard.Dealer {
			break
		}
	}
	pokerBoard.GameState = "preFlop"
	pokerBoard.Dealer.Next_player.CurrentBet = 10
	pokerBoard.Dealer.Next_player.Money -= 10
	pokerBoard.Dealer.Next_player.Next_player.CurrentBet = 20
	pokerBoard.Dealer.Next_player.Next_player.Money -= 20
	pokerBoard.CurrentBet = 20
	pokerBoard.Pot = 30
	pokerBoard.Starter = pokerBoard.Dealer.Next_player.Next_player.Next_player
	pokerBoard.CurrentPlayer = pokerBoard.Starter

}

func goFlopStuff(pokerBoard *board.Board) {
	var i = pokerBoard.Dealer
	var card1 = pokerBoard.Deck.GetPokerCard()
	var card2 = pokerBoard.Deck.GetPokerCard()
	var card3 = pokerBoard.Deck.GetPokerCard()
	for {
		sendPokerMessage(card1, i.Conn)
		sendPokerMessage(card2, i.Conn)
		sendPokerMessage(card3, i.Conn)
		i = i.Next_player
		i.CurrentBet = 0
		if i == pokerBoard.Dealer {
			break
		}
	}
	pokerBoard.CurrentBet = 0
	pokerBoard.GameState = "afterFlop"
	pokerBoard.Starter = pokerBoard.Dealer.Next_player
}

func goTurnStuff(pokerBoard *board.Board) {
	var i = pokerBoard.Dealer
	var card4 = pokerBoard.Deck.GetPokerCard()
	for {
		sendPokerMessage(card4, i.Conn)
		i = i.Next_player
		i.CurrentBet = 0
		if i == pokerBoard.Dealer {
			break
		}
	}
	pokerBoard.CurrentBet = 0
	pokerBoard.GameState = "afterTurn"
	pokerBoard.Starter = pokerBoard.Dealer.Next_player
}

func goRiverStuff(pokerBoard *board.Board) {
	var i = pokerBoard.Dealer
	var card5 = pokerBoard.Deck.GetPokerCard()
	for {
		sendPokerMessage(card5, i.Conn)
		i = i.Next_player
		i.CurrentBet = 0
		if i == pokerBoard.Dealer {
			break
		}
	}
	pokerBoard.CurrentBet = 0
	pokerBoard.GameState = "afterRiver"
	pokerBoard.Starter = pokerBoard.Dealer.Next_player
}
