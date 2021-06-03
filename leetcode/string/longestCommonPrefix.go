package string

// 编写一个函数来查找字符串数组中的最长公共前缀。
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result = commonPrefix(result, strs[i])
		if result == "" {
			break
		}
	}
	return result
}
func commonPrefix(str1, str2 string) string {
	if len(str2) < len(str1) {
		return commonPrefix(str2, str1)
	}
	// str1, str2 = str2, str1
	var commonLength int
	for commonLength < len(str1) {
		if str1[commonLength] == str2[commonLength] {
			commonLength++
		} else {
			break
		}
	}
	return str1[0:commonLength]
}
