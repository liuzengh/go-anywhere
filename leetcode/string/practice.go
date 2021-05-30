package string

// 给定一个字符串，请你找出其中不含有重复字符的 最长子串 的长度。
func lengthOfLongestSubstring(s string) int {
	hasSeen := make(map[byte]int)
	var result int
	for i, j := 0, 0; i < len(s); i++ {
		if _, exist := hasSeen[s[i]]; exist && hasSeen[s[i]] >= j {
			j = hasSeen[s[i]] + 1
		}
		cur_len := i - j + 1
		if cur_len > result {
			result = cur_len
		}
		hasSeen[s[i]] = i
	}
	return result
}
