func singleNumber(nums []int) int {
	count := make(map[int]int)

	for _, num := range nums {
		count[num]++
	}

	for _, num := range nums {
		if count[num] == 1 {
			return num
		}
	}

	return -1
}