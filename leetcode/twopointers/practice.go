package twopointers

// 给你 n 个非负整数 a1，a2，...，an，每个数代表坐标中的一个点 (i, ai) 。
// 在坐标内画 n 条垂直线，垂直线 i 的两个端点分别为 (i, ai) 和 (i, 0) 。
// 找出其中的两条线，使得它们与 x 轴共同构成的容器可以容纳最多的水。
func maxArea(height []int) int {
	left, right := 0, len(height)-1
	var volume int
	for left <= right {
		temp := (right - left) * MinInt(height[left], height[right])
		if temp > volume {
			volume = temp
		}
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}
	return volume
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
