func timeRequiredToBuy(tickets []int, k int) int {
	queue := make([]int, len(tickets))

	for i := 0; i < len(tickets); i++ {
		queue[i] = i
	}

	time := 0

	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]
		time++
		tickets[front]--

		if front == k && tickets[front] == 0 {
			return time
		}

		if tickets[front] == 0 {
			continue
		}

		queue = append(queue, front)
	}

	return time
}