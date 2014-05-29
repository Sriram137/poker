package board

import (
	"github.com/elricL/poker/cards"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"strings"
)

type Player struct {
	Next_player *Player
	Folded      bool
	Conn        *websocket.Conn
	Name        string
	Hand        []string
	CurrentBet  int
	Money       int
}

type Board struct {
	Deck          cards.Deck
	Dealer        *Player
	Starter       *Player
	CurrentPlayer *Player
	GameState     string
	BoardCards    []string
	CurrentBet    int
	Pot           int
}

func (b *Board) String() string {
	return strings.Join([]string{"BOARD", "Dealer", b.Dealer.Name, "Starter", b.Starter.Name, "Current", b.CurrentPlayer.Name, b.GameState, strconv.Itoa(b.CurrentBet), strconv.Itoa(b.Pot)}, " ")
}

func (P *Player) String() string {
	return strings.Join([]string{"PLAYER", P.Name, strconv.Itoa(P.CurrentBet), strconv.Itoa(P.Money)}, " ")
}

func (P *Player) PlayerInfo() string {
	var k string = ""
	if P.Folded == false {
		k = "not"
	}
	var hand = P.HandInfo()
	return strings.Join([]string{"You have", k, "folded.\nYour Hand :", hand, "\nCurrent Bet :", strconv.Itoa(P.CurrentBet), "\nMoney Left : ", strconv.Itoa(P.Money)}, " ")
}

func (P *Player) HandInfo() string {
	if P.Hand == nil {
		return "empty"
	}
	return strings.Join(P.Hand, " ")
}

func (b *Board) PrintCards() string {
	return strings.Join(b.BoardCards, " ")
}

func MakeNewBoard() Board {
	return Board{cards.Deck{}, nil, nil, nil, "waiting", make([]string, 0), 0, 0}
}

func (b *Board) Shuffle() {
	b.Deck = cards.Deck{}
	b.Deck.MakeShuffledCardPack()
}

func (board *Board) AddPlayer(player Player) {
	if board.Dealer == nil {
		board.Dealer = &player
		player.Next_player = board.Dealer
		log.Println("start")
		return
	}
	if board.Dealer == board.Dealer.Next_player {
		board.Dealer.Next_player = &player
		player.Next_player = board.Dealer
		log.Println("second")
		return
	}
	var start = board.Dealer.Next_player
	for ; start.Next_player != board.Dealer; start = start.Next_player {
	}
	start.Next_player = &player
	player.Next_player = board.Dealer
	log.Println("third")
}

func (board *Board) Length() int {
	if board.Dealer == nil {
		return 0
	}
	if board.Dealer == board.Dealer.Next_player {
		return 1
	}
	var count = 0
	for start := board.Dealer.Next_player; start != board.Dealer; start = start.Next_player {
		count++
	}
	count++
	return count
}

func (board *Board) Print() {
	if board.Dealer == nil {
		log.Println("first_p")
		return
	}
	if board.Dealer == board.Dealer.Next_player {
		log.Println("second_p")
		return
	}
	log.Println("thrid_p")
	var start = board.Dealer.Next_player
	for ; start != board.Dealer; start = start.Next_player {
		log.Println(*start)
	}
	log.Println(*start)
}

func (pokerPlayer *Player) FindNextUnfoldedPlayer() *Player {
	var i = pokerPlayer.Next_player
	for {
		if i.Folded {
			i = i.Next_player
		} else {
			return i
		}
	}
}
