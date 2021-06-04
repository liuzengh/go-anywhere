package main

import (
	"sort"
)

// 给你一个包含 n 个整数的数组 nums，判断 nums 中是否存在三个元素 a，b，c ，使得 a + b + c = 0 ？
// 请你找出所有和为 0 且不重复的三元组。
func threeSum(nums []int) [][]int {
	result := make([][]int, 0)
	sort.Ints(nums)
	for first := 0; first < len(nums)-2; first++ {
		if 0 < first && nums[first-1] == nums[first] {
			continue
		}
		others := -nums[first]
		third := len(nums) - 1
		for second := first + 1; second < third; second++ {
			if first+1 < second && nums[second-1] == nums[second] {
				continue
			}
			for second < third && others < nums[second]+nums[third] {
				third--
			}

			if second < third && nums[second]+nums[third] == others {
				result = append(result, []int{nums[first], nums[second], nums[third]})
			}
		}
	}
	return result
}

// 给定一个包括 n 个整数的数组 nums 和 一个目标值 target。
// 找出 nums 中的三个整数，使得它们的和与 target 最接近。
// 返回这三个数的和。假定每组输入只存在唯一答案。
func threeSumClosest(nums []int, target int) int {
	result := 100000
	sort.Ints(nums)
	for first := 0; first < len(nums)-2; first++ {
		if 0 < first && nums[first-1] == nums[first] {
			continue
		}
		others := target - nums[first]
		third := len(nums) - 1
		for second := first + 1; second < third; second++ {
			if first+1 < second && nums[second-1] == nums[second] {
				continue
			}
			for second < third && others < nums[second]+nums[third] {
				if absInt(others-nums[second]-nums[third]) < absInt(target-result) {
					result = nums[first] + nums[second] + nums[third]
				}
				third--
			}

			if second < third {
				if absInt(others-nums[second]-nums[third]) < absInt(target-result) {
					result = nums[first] + nums[second] + nums[third]
				}
			}
		}
	}
	return result
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
