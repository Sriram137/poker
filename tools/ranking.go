package tools

func rank(cards []int) ([]int, []int) {
	freqs := make([]int, 13)
	out1 := make([]int, 5)
	out2 := make([]int, 5)
	i1 := 0
	i2 := 0
	for _, v := range cards {
		freqs[12-v+2]++
	}
	for i, v := range freqs {
		if v == 4 {
			out1[i1] = v
			out2[i2] = 14 - i
			i1++
			i2++
		}
	}
	for i, v := range freqs {
		if v == 3 {
			out1[i1] = v
			out2[i2] = 14 - i
			i1++
			i2++
		}
	}
	for i, v := range freqs {
		if v == 2 {
			out1[i1] = v
			out2[i2] = 14 - i
			i1++
			i2++
		}
	}
	for i, v := range freqs {
		if v == 1 {
			out1[i1] = v
			out2[i2] = 14 - i
			i1++
			i2++
		}
	}
	return out1, out2
}

func findRank(cards []int) ([]int, int) {
	maxHand := make([]int, 5)
	maxRank := 0
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 7; j++ {
			temp1 := make([]int, i)
			temp2 := make([]int, j-i-1)
			temp3 := make([]int, 6-j)
			copy(temp1, cards[:i])
			copy(temp2, cards[i+1:j])
			copy(temp3, cards[j+1:])
			curHand := append(append(temp1, temp2...), temp3...)
			curRank := dummy(curHand)
			if curRank > maxRank {
				maxRank = curRank
				maxHand = curHand
			}
		}
	}
	return maxHand, maxRank
}

func dummy(cards []int) int {
	return cards[4]
}
