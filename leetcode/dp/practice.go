package dp

type interval struct {
	X int
	Y int
}

func longestPalindrome(s string) string {
	if len(s) == 0 || len(s) == 1 {
		return s
	}
	is_palindrome := make([][]bool, len(s))
	for row := range is_palindrome {
		is_palindrome[row] = make([]bool, len(s))
	}
	var result interval
	for i := 0; i < len(s); i++ {
		is_palindrome[i][i] = true
	}
	for cnt := 2; cnt <= len(s); cnt++ {
		for x := 0; x <= len(s)-cnt; x++ {
			y := x + cnt - 1
			if s[x] == s[y] && (y-1 < x+1 || is_palindrome[x+1][y-1]) {
				result = interval{x, y}
				is_palindrome[x][y] = true
			}
		}
	}
	return s[result.X : result.Y+1]
}

type Pos struct {
	X, Y int
}

// 给你一个字符串 s 和一个字符规律 p，请你来实现一个支持 '.' 和 '*' 的正则表达式匹配。
func isMatch(s string, p string) bool {
	memo := make(map[Pos]bool)
	var dp func(x, y int) bool
	dp = func(x, y int) bool {
		pos := Pos{X: x, Y: y}
		var result bool
		if memo[pos] {
			return memo[pos]
		}
		if y == len(p) {
			result = (x == len(s))
		} else {
			first_match := x < len(s) && (p[y] == s[x] || p[y] == '.')
			if y+1 < len(p) && p[y+1] == '*' {
				result = (dp(x, y+2)) || (first_match && dp(x+1, y))
			} else {
				result = first_match && dp(x+1, y+1)
			}
		}
		memo[pos] = result
		return result
	}

	return dp(0, 0)
}
