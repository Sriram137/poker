package ranking

import (
	"github.com/elricL/poker/board"
	"log"
	"reflect"
	"sort"
)

func card_ranks(hand []string) []int {
	var ranks []int = make([]int, 6)
	rankSpecial := []int{14, 5, 4, 3, 2}
	for i, v := range hand {
		if v[0] == 'A' {
			ranks[i] = 14
		} else if v[0] == 'K' {
			ranks[i] = 13
		} else if v[0] == 'Q' {
			ranks[i] = 12
		} else if v[0] == 'J' {
			ranks[i] = 11
		} else if v[0] == 'T' {
			ranks[i] = 10
		} else {
			ranks[i] = int(v[0] - '0')
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ranks[0:5])))
	if reflect.DeepEqual(ranks, rankSpecial) {
		ranks = []int{5, 4, 3, 2, 1}
	}
	return ranks
}

func hand_rank(hand []string) []int {
	ranks := card_ranks(hand)
	var ret []int
	if straight(ranks) && flush(hand) {
		ret = append(ret, 8, max(ranks))
		return ret
	} else if kind(4, ranks) != -1 {
		ret = append(ret, 7, kind(4, ranks), kind(1, ranks))
		return ret
	} else if kind(3, ranks) != -1 && kind(2, ranks) != -1 {
		ret = append(ret, 6, kind(3, ranks), kind(2, ranks))
		return ret
	} else if flush(hand) {
		ret = append(ret, 5)
		ret = append(ret, ranks[0:5]...)
		return ret
	} else if straight(ranks) {
		ret = append(ret, 4, max(ranks))
		return ret
	} else if kind(3, ranks) != -1 {
		ret = append(ret, 3, kind(3, ranks))
		ret = append(ret, ranks[0:5]...)
		return ret
	} else if a, b := two_pairs(ranks); a != -1 && b != -1 {
		a, b = two_pairs(ranks)
		ret = append(ret, 2, a, b)
		ret = append(ret, ranks[0:5]...)
		return ret
	} else if kind(2, ranks) != -1 {
		ret = append(ret, 1, kind(2, ranks))
		ret = append(ret, ranks[0:5]...)
		return ret
	} else {
		ret = append(ret, 0)
		ret = append(ret, ranks[0:5]...)
		return ret
	}
}

func max(ranks []int) int {
	m := 0
	for i, v := range ranks {
		if i == 0 {
			m = v
		} else {
			if v > m {
				m = v
			}
		}
	}
	return m
}

func straight(ranks []int) bool {
	len := len(ranks)
	if ranks[0]-ranks[4] == 4 && len == 5 {
		return true
	} else {
		return false
	}
}

func two_pairs(ranks_orig []int) (int, int) {
	ranks := make([]int, 10)
	copy(ranks, ranks_orig)
	pair := kind(2, ranks)
	sort.Sort(sort.IntSlice(ranks[0:5]))
	lowpair := kind(2, ranks)
	if pair != -1 && pair != lowpair {
		return pair, lowpair
	} else {
		return -1, -1
	}
}

func flush(hand []string) bool {
	var suit string
	count := 0
	for i, v := range hand {
		if i == 0 {
			suit = string(v[1])
			count = 1
		} else {
			if string(v[1]) == suit {
				count++
			}
		}
	}
	if count == 5 {
		return true
	} else {
		return false
	}
}

func kind(n int, ranks []int) (rank int) {
	for _, v := range ranks {
		if count(v, ranks) == n {
			return v
		}
	}
	return -1
}

func count(r int, ranks []int) int {
	count := 0
	for _, v := range ranks {
		if v == r {
			count++
		}
	}
	return count
}

func FindWinners(pokerBoard *board.Board) []*board.Player {
	var winners []*board.Player
	bestHand := make([]int, 10)
	for i := 0; i < 10; i++ {
		bestHand[i] = 0
	}
	firstPlayer := pokerBoard.Dealer
	curPlayer := firstPlayer
	for {
		if curPlayer.Folded {
			continue
		}
		curHand := findBestHand(append(curPlayer.Hand, pokerBoard.BoardCards...))

		log.Println(curHand)

		if compareHands(curHand, bestHand) > 0 {
			log.Println("HELLO1")
			bestHand = curHand
			winners = make([]*board.Player, 0)
			winners = append(winners, curPlayer)
			log.Println(winners)
		} else if compareHands(curHand, bestHand) == 0 {
			log.Println("HELLO2")
			winners = append(winners, curPlayer)
		}
		curPlayer = curPlayer.Next_player
		if curPlayer == firstPlayer {
			break
		}
	}
	log.Println(winners)
	return winners
}

func findBestHand(cards []string) []int {
	maxHand := make([]int, 10)
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 7; j++ {
			temp1 := make([]string, i)
			temp2 := make([]string, j-i-1)
			temp3 := make([]string, 6-j)
			copy(temp1, cards[:i])
			copy(temp2, cards[i+1:j])
			copy(temp3, cards[j+1:])
			curHand := hand_rank(append(append(temp1, temp2...), temp3...))
			if compareHands(curHand, maxHand) > 0 {
				maxHand = curHand
			}
		}
	}
	return maxHand
}

func compareHands(firstHand []int, secondHand []int) int {
	for i, v := range firstHand {
		if v > secondHand[i] {
			return 1
		}
		if v < secondHand[i] {
			return -1
		}
	}
	return 0
}

// func test() {
// h := []string{"2S", "2D", "2C", "2E", "5S"}
// ret := hand_rank(h)
// }
