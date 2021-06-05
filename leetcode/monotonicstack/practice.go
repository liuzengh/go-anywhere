package monotonicstack

import "math"

func trap(height []int) int {
	s, top := []int{}, -1
	var result int
	for index, item := range height {
		if top == -1 || item <= height[s[top]] {
			top++
			s = append(s, index)
		} else {
			for -1 < top && height[s[top]] < item {
				top--
				if top == -1 {
					break
				}
				distance := index - s[top] - 1
				delta := MinInt(item, height[s[top]]) - height[s[top+1]]
				result += distance * delta

			}
			s = s[0 : top+1]
			s = append(s, index)
			top++
		}
	}
	return result
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

/*
func nextGreaterElement(nums1 []int, nums2 []int) []int {
	result := make([]int, len(nums1))
	for index1, item1 := range nums1 {
		result[index1] = -1
		targetIndex := len(nums2)
		for index2, item2 := range nums2 {
			if item1 == item2 {
				targetIndex = index2
			}
			if item1 < item2 && targetIndex < index2 {
				result[index1] = item2
				break
			}
		}
	}
	return result
}
*/

func nextGreaterElement(nums1 []int, nums2 []int) []int {
	result := make([]int, len(nums1))
	nextGreater := make([]int, len(nums2))
	stack := []int{math.MinInt32}
	elem2index := map[int]int{}
	for i := len(nums2) - 1; 0 <= i; i-- {
		for 0 < len(stack) && stack[len(stack)-1] <= nums2[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			nextGreater[i] = -1
		} else {
			nextGreater[i] = stack[len(stack)-1]
		}
		stack = append(stack, nums2[i])
		elem2index[nums2[i]] = i
	}
	for index1, item1 := range nums1 {
		result[index1] = nextGreater[elem2index[item1]]
	}
	return result
}

func dailyTemperatures(temperatures []int) []int {
	stack := []int{}
	result := make([]int, len(temperatures))
	for i := len(temperatures) - 1; 0 <= i; i-- {
		for 0 < len(stack) && temperatures[stack[len(stack)-1]] <= temperatures[i] {
			stack = stack[:len(stack)-1]
		}
		if 0 < len(stack) {
			result[i] = stack[len(stack)-1] - i
		}
		stack = append(stack, i)
	}
	return result
}

func nextGreaterElements(nums []int) []int {
	nums = append(nums, nums...)
	result := make([]int, len(nums))
	stack := []int{}

	for i := len(nums) - 1; 0 <= i; i-- {
		for 0 < len(stack) && stack[len(stack)-1] <= nums[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			result[i] = -1
		} else {
			result[i] = stack[len(stack)-1]
		}
		stack = append(stack, nums[i])
	}
	return result[:len(result)/2]
}

func largestRectangleArea(heights []int) int {
	left, right := make([]int, len(heights)), make([]int, len(heights))
	stack := []int{}

	for i := 0; i < len(heights); i++ {
		for 0 < len(stack) && heights[i] < heights[stack[len(stack)-1]] {
			curPos := stack[len(stack)-1]
			right[curPos] = (i - curPos) * heights[curPos]
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	for _, curPos := range stack {
		right[curPos] = (stack[len(stack)-1] + 1 - curPos) * heights[curPos]
	}
	// fmt.Println(right)
	stack = []int{}
	for i := len(heights) - 1; 0 <= i; i-- {
		for 0 < len(stack) && heights[i] < heights[stack[len(stack)-1]] {
			curPos := stack[len(stack)-1]
			left[curPos] = (curPos - i) * heights[curPos]
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	for _, curPos := range stack {
		left[curPos] = (curPos - stack[len(stack)-1] + 1) * heights[curPos]
	}
	// fmt.Println(left)
	result := 0
	for index, height := range heights {
		area := right[index] + left[index] - height
		if result < area {
			result = area
		}
	}
	return result

}
