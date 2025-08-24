func isValid(s string) bool {
	left := []rune{}

	for _, c := range s {
		if c == '(' || c == '[' || c == '{' {
			left = append(left, c)
		} else {
			if len(left) > 0 && leftOf(c) == left[len(left)-1] {
				left = left[:len(left)-1]
			} else {
				return false
			}
		}
	}

	return len(left) == 0
}

func leftOf(c rune) rune {
	if c == '}' {
		return '{'
	} else if c == ')' {
		return '('
	} else {
		return '['
	}
}