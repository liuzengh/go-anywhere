package math

import "math"

/*
func reverse(x int) int {
	digits := make([]int, 0, 32)
	sign := true
	if x == 0 {
		return x
	} else if x < 0 {
		sign = false
		x = -x
	}
	for x > 0 {
		digits = append(digits, x%10)
		x /= 10
	}
	var result int
	for i := 0; i < len(digits); i++ {
		if  i == 9 {
			// first9 := (2**31) / 10
            first9 := math.MaxInt32 / 10
			if sign {
				if result > first9 || (result==first9 && digits[i] > 7) {
					return 0
				}
			} else {
				if result > first9 || (result==first9 && digits[i] > 8) {
					return 0
				}
			}

		}
		result = result * 10 + digits[i]
	}
	if !sign {
		result = -result
	}
	return result
}
*/
// 给你一个 32 位的有符号整数 x ，返回将 x 中的数字部分反转后的结果。
// 如果反转后整数超过 32 位的有符号整数的范围 [−2**31,  2**31 − 1] ，就返回 0。

func reverse(x int) int {
	var result int
	first9 := math.MaxInt32 / 10
	last9 := math.MinInt32 / 10
	for x != 0 {
		if result > first9 || (result == first9 && x%10 > 7) {
			return 0
		}
		if result < last9 || (result == last9 && x%10 < -8) {
			return 0
		}
		result = result*10 + x%10
		x /= 10
	}
	return result
}

// 给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	digits := make([]int, 0, 32)
	// digits := []int{}
	for x > 0 {
		digits = append(digits, x%10)
		x /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		if digits[i] != digits[j] {
			return false
		}
	}
	return true
}
