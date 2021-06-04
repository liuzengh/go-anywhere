package main

func letterCombinations(digits string) []string {
	button := map[byte]string{
		'2': "abc",
		'3': "def",
		'4': "ghi",
		'5': "jkl",
		'6': "mno",
		'7': "pqrs",
		'8': "tuv",
		'9': "wxyz",
	}
	temp := make([]byte, len(digits))

	var result []string
	search(0, digits, temp, &result, button)
	return result
}

func search(cur int, digits string, temp []byte, result *[]string, button map[byte]string) {
	if cur == len(digits) {
		*result = append(*result, string(temp))
		return
	}
	for _, ch := range button[digits[cur]] {
		temp[cur] = byte(ch)
		search(cur+1, digits, temp, result, button)
	}
}
