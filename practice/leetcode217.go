func containsDuplicate(nums []int) bool {
	count := make(map[int]bool)

	for _, num := range nums {
		if count[num] {
			return true
		}

		count[num] = true
	}

	return false
}