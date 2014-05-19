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
