package stack

func isValid(s string) bool {
	if len(s)%2 == 1 {
		return false
	}
	stack := make([]rune, 0, len(s))
	leftMatch := map[rune]rune{
		'(': ')',
		'{': '}',
		'[': ']',
	}
	rightMatch := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}
	for _, ch := range s {
		if _, exist := leftMatch[ch]; exist {
			stack = append(stack, ch)
		} else {
			if len(stack) != 0 && stack[len(stack)-1] == rightMatch[ch] {
				stack = stack[:len(stack)-1]
			} else {
				return false
			}
		}
	}
	return len(stack) == 0
}
