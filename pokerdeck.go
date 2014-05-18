package main

import (
	"math/rand"
)

func makeShuffledCardPack() {
	curStart = -1
	var niceDeck = make([]string, 52)
	var count = 0
	for _, num := range "123456789TJKQ" {
		for _, s := range "DHSC" {
			niceDeck[count] = string(num) + string(s)
			count += 1
		}
	}
	var randS = rand.Perm(51)
	for i, j := range randS {
		cardDeck[i] = niceDeck[j]
	}
}

func getPokerCard() string {
	curStart += 1
	return cardDeck[curStart]
}
