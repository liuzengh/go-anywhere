package string

import "math"

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

// 将一个给定字符串 s 根据给定的行数 numRows ，以从上往下、从左到右进行 Z 字形排列。
func convert(s string, numRows int) string {
	if numRows == 1 || len(s) <= numRows {
		return s
	}
	display := make([][]byte, numRows)
	x, down := 0, true
	for i := 0; i < len(s); i++ {
		display[x] = append(display[x], s[i])
		if down {
			x++
			if x == numRows {
				down = false
				x -= 2
			}
		} else {
			x--
			if x == -1 {
				down = true
				x += 2
			}
		}
	}
	result := make([]byte, 0, len(s))
	for _, row := range display {
		for _, item := range row {
			result = append(result, item)
		}
	}
	return string(result)
}

// 请你来实现一个 myAtoi(string s) 函数，使其能将字符串转换成一个 32 位有符号整数（类似 C/C++ 中的 atoi 函数）
func myAtoi(s string) int {
	state := 1
	sign := true
	var result int
	for _, ch := range s {
		if state == 1 {
			if ch == ' ' {

			} else if ch == '+' || ch == '-' {
				state = 3
				if ch == '-' {
					sign = false
				}
			} else if '0' <= ch && ch <= '9' {
				state = 4
				result = int(ch - '0')
			} else {
				break
			}
		} else if state == 3 {
			if '0' <= ch && ch <= '9' {
				state = 4
				result = int(ch - '0')
			} else {
				break
			}
		} else {
			// state == 4
			if ch < '0' || '9' < ch {
				break
			}
			result = result*10 + int(ch-'0')
			if sign && result > math.MaxInt32 {
				return math.MaxInt32
			}
			if !sign && result > math.MaxInt32+1 {
				return math.MinInt32
			}

		}
	}
	// state 2

	if !sign {
		result = -result
	}
	return result
}
