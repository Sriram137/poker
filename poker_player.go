package main

import (
	"github.com/gorilla/websocket"
)

type orderedSet struct {
	items    []string
	itemsMap map[string]int
	posMap   map[string]int
	connMap  map[string]*websocket.Conn
	cardMap  map[string][]string
}

type Player struct {
	next_player *Player
	folded      bool
	conn        *websocket.Conn
	name        string
	hand        []string
}

type Board struct {
	cardDeck []string
	dealer   *Player
	size     int
}

func (board *Board) addPlayer(player Player) {
	if board.dealer == nil {
		*board.dealer = player
		player.next_player = board.dealer
		return
	}
	if board.dealer == board.dealer.next_player {
		*board.dealer.next_player = player
		player.next_player = board.dealer
		return
	}
	var start = board.dealer.next_player
	for ; start.next_player != board.dealer; start = start.next_player {
	}
	*start.next_player = player
	player.next_player = board.dealer
}

func (o *orderedSet) addPlayer(name string, conn *websocket.Conn) {
	if o.itemsMap[name] == 0 {
		o.itemsMap[name] += 1
		o.posMap[name] = len(o.items)
		o.connMap[name] = conn
		o.items = append(o.items, name)
	}
}

func (o *orderedSet) getPostion(name string) int {
	return o.posMap[name]
}

func (o *orderedSet) length() int {
	return len(o.items)
}
