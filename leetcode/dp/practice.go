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
