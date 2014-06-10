package game

import (
	"github.com/elricL/poker/board"
	"github.com/elricL/poker/ranking"
	"log"
	"strconv"
)

func gameStart(pokerBoard *board.Board) {
	sendAll(pokerBoard, "Commencing PREFLOP")
	pokerBoard.BoardCards = []string{"__", "__", "__", "__", "__"}
	pokerBoard.Shuffle()
	pokerBoard.Dealer = pokerBoard.Dealer.Next_player
	var i = pokerBoard.Dealer
	var count = 0
	for {
		var card1 = pokerBoard.Deck.GetPokerCard()
		var card2 = pokerBoard.Deck.GetPokerCard()
		sendPokerMessage(card1, i.Conn)
		sendPokerMessage(card2, i.Conn)
		sendPokerMessage("Your postition "+strconv.Itoa(count), i.Conn)
		i.Hand = []string{card1, card2}
		i = i.Next_player
		i.Folded = false
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
	sendAll(pokerBoard, getPlayerPositionsInPrintableFormat(pokerBoard))

}

func goFlopStuff(pokerBoard *board.Board) {
	var i = pokerBoard.Dealer
	var card1 = pokerBoard.Deck.GetPokerCard()
	var card2 = pokerBoard.Deck.GetPokerCard()
	var card3 = pokerBoard.Deck.GetPokerCard()
	pokerBoard.BoardCards = []string{card1, card2, card3, "__", "__"}
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
	pokerBoard.Starter = pokerBoard.Dealer.FindNextUnfoldedPlayer()
	pokerBoard.CurrentPlayer = pokerBoard.Dealer.FindNextUnfoldedPlayer()
}

func goTurnStuff(pokerBoard *board.Board) {
	var i = pokerBoard.Dealer
	var card4 = pokerBoard.Deck.GetPokerCard()
	pokerBoard.BoardCards[3] = card4
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
	pokerBoard.Starter = pokerBoard.Dealer.FindNextUnfoldedPlayer()
	pokerBoard.CurrentPlayer = pokerBoard.Dealer.FindNextUnfoldedPlayer()
}

func goRiverStuff(pokerBoard *board.Board) {
	var i = pokerBoard.Dealer
	var card5 = pokerBoard.Deck.GetPokerCard()
	pokerBoard.BoardCards[4] = card5
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
	pokerBoard.Starter = pokerBoard.Dealer.FindNextUnfoldedPlayer()
	pokerBoard.CurrentPlayer = pokerBoard.Dealer.FindNextUnfoldedPlayer()
}

func findNextUnfoldedPlayer(pokerPlayer *board.Player) *board.Player {
	var i = pokerPlayer.Next_player
	for {
		if i.Folded {
			i = i.Next_player
		} else {
			return i
		}
	}
}

func findGameWinner(pokerBoard *board.Board) ([]*board.Player, int) {
	winners := make([]*board.Player, 0)
	starter := pokerBoard.Starter
	if starter.FindNextUnfoldedPlayer() == starter {
		starter.Money += pokerBoard.Pot
		log.Println("WTF")
		log.Println(starter)
		log.Println(starter.FindNextUnfoldedPlayer())
		winners := append(winners, starter)
		amount := pokerBoard.Pot / len(winners)
		log.Println(winners)
		log.Println(amount)
		return winners, amount
	}
	winners = ranking.FindWinners(pokerBoard)
	amount := pokerBoard.Pot / len(winners)
	for _, winner := range winners {
		winner.Money += amount
	}
	return winners, amount
}

func resetGame(pokerBoard *board.Board) {
	pokerBoard.Deck.MakeShuffledCardPack()
	pokerBoard.CurrentBet = 0
	pokerBoard.Pot = 0
	var i = pokerBoard.Dealer
	for {
		i.CurrentBet = 0
		i.Folded = false
		if i == pokerBoard.Dealer {
			break
		}
	}
	gameStart(pokerBoard)
}
