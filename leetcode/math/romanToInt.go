package math

// 给定一个罗马数字，将其转换成整数。输入确保在 1 到 3999 的范围内。
func romanToInt(s string) int {
	// vals := [...]int {1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	// syms := [...]string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	var result int
	for i := 0; i < len(s); i++ {
		if s[i] == 'M' {
			result += 1000
		} else if s[i] == 'C' {
			if i+1 < len(s) {
				if s[i+1] == 'M' {
					i++
					result += 900
				} else if s[i+1] == 'D' {
					i++
					result += 400
				} else {
					result += 100
				}
			} else {
				result += 100
			}
		} else if s[i] == 'D' {
			result += 500
		} else if s[i] == 'X' {
			if i+1 < len(s) {
				if s[i+1] == 'C' {
					i++
					result += 90
				} else if s[i+1] == 'L' {
					i++
					result += 40
				} else {
					result += 10
				}
			} else {
				result += 10
			}
		} else if s[i] == 'L' {
			result += 50
		} else if s[i] == 'I' {
			if i+1 < len(s) {
				if s[i+1] == 'X' {
					i++
					result += 9
				} else if s[i+1] == 'V' {
					i++
					result += 4
				} else {
					result += 1
				}
			} else {
				result += 1
			}
		} else if s[i] == 'V' {
			result += 5
		}
	}
	return result
}
