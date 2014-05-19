package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	next_player *Player
	folded      bool
	conn        *websocket.Conn
	name        string
	hand        []string
}

type Board struct {
	deck      Deck
	dealer    *Player
	starter   *Player
	gameState string
}

func makeNewBoard() {
	board = Board{Deck{}, nil, nil, "waiting"}
}

func (board *Board) addPlayer(player Player) {
	if board.dealer == nil {
		board.dealer = &player
		player.next_player = board.dealer
		log.Println("start")
		log.Println(board)
		return
	}
	if board.dealer == board.dealer.next_player {
		board.dealer.next_player = &player
		player.next_player = board.dealer
		log.Println("second")
		log.Println(board)
		return
	}
	var start = board.dealer.next_player
	for ; start.next_player != board.dealer; start = start.next_player {
	}
	start.next_player = &player
	player.next_player = board.dealer
	log.Println("third")
	log.Println(board)
}

func (board *Board) length() int {
	if board.dealer == nil {
		return 0
	}
	if board.dealer == board.dealer.next_player {
		return 1
	}
	var count = 0
	for start := board.dealer.next_player; start != board.dealer; start = start.next_player {
		count++
	}
	count++
	return count
}

func (board *Board) print() {
	if board.dealer == nil {
		log.Println("first_p")
		log.Println(*board.dealer)
		return
	}
	if board.dealer == board.dealer.next_player {
		log.Println("second_p")
		log.Println(*board.dealer)
		log.Println(*board.dealer.next_player)
		return
	}
	log.Println("thrid_p")
	var start = board.dealer.next_player
	for ; start != board.dealer; start = start.next_player {
		log.Println(*start)
	}
	log.Println(*start)
}
