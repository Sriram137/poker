package main

import (
	"github.com/gorilla/websocket"
)

var currentPlayers = orderedSet{make([]string, 0), make(map[string]int), make(map[string]int), make(map[string]*websocket.Conn), make(map[string][]string)}
var waitingPlayers = orderedSet{make([]string, 0), make(map[string]int), make(map[string]int), make(map[string]*websocket.Conn), make(map[string][]string)}

var board = Board{cardDeck, nil, 0}

var inProgress = false

var cardDeck = make([]string, 52)
var curStart = -1

var gameState = "waiting"
var dealer = 0

var expectedPlayer = 0
