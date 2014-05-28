package tools

func rank(cards []int) ([]int, []int) {
	freqs := make([]int, 13)	
	out1 := make([]int, 5)	
	out2 := make([]int, 5)	
	i1 := 0
	i2 := 0
	for _,v := range cards {
		freqs[12-v+2]++
	}
	for i,v := range freqs {
		if v == 4 {
			out1[i1]=v
			out2[i2]=14-i
			i1++
			i2++
		}
	}
	for i,v := range freqs {
		if v == 3 {
			out1[i1]=v
			out2[i2]=14-i
			i1++
			i2++
		}
	}
	for i,v := range freqs {
		if v == 2 {
			out1[i1]=v
			out2[i2]=14-i
			i1++
			i2++
		}
	}
	for i,v := range freqs {
		if v == 1 {
			out1[i1]=v
			out2[i2]=14-i
			i1++
			i2++
		}
	}
	return out1, out2
}
