package cards

import (
	"math/rand"
	"time"
)

type Deck struct {
	cardDeck []string
	iterator int
}

func (deck *Deck) MakeShuffledCardPack() {
	deck.iterator = -1
	var niceDeck = make([]string, 52)
	deck.cardDeck = make([]string, 52)

	var count = 0
	for _, num := range "123456789TJKQ" {
		for _, s := range "DHSC" {
			niceDeck[count] = string(num) + string(s)
			count += 1
		}
	}
	rand.Seed(time.Now().Unix())
	var randS = rand.Perm(51)
	for i, j := range randS {
		deck.cardDeck[i] = niceDeck[j]
	}
}

func (deck *Deck) GetPokerCard() string {
	deck.iterator += 1
	return deck.cardDeck[deck.iterator]
}
