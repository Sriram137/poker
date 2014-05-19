package board

import (
	"github.com/elricL/poker/cards"
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	Next_player *Player
	Folded      bool
	Conn        *websocket.Conn
	Name        string
	Hand        []string
}

type Board struct {
	Deck      cards.Deck
	Dealer    *Player
	Starter   *Player
	GameState string
}

func MakeNewBoard() Board {
	return Board{cards.Deck{}, nil, nil, "waiting"}
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
		log.Println(board)
		return
	}
	if board.Dealer == board.Dealer.Next_player {
		board.Dealer.Next_player = &player
		player.Next_player = board.Dealer
		log.Println("second")
		log.Println(board)
		return
	}
	var start = board.Dealer.Next_player
	for ; start.Next_player != board.Dealer; start = start.Next_player {
	}
	start.Next_player = &player
	player.Next_player = board.Dealer
	log.Println("third")
	log.Println(board)
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
		log.Println(*board.Dealer)
		return
	}
	if board.Dealer == board.Dealer.Next_player {
		log.Println("second_p")
		log.Println(*board.Dealer)
		log.Println(*board.Dealer.Next_player)
		return
	}
	log.Println("thrid_p")
	var start = board.Dealer.Next_player
	for ; start != board.Dealer; start = start.Next_player {
		log.Println(*start)
	}
	log.Println(*start)
}
