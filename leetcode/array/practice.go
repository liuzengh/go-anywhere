package array

// 给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
func twoSum(nums []int, target int) []int {
	haveSeen := make(map[int]int)
	for index_one, one := range nums {
		another := target - one
		if _, exist := haveSeen[another]; exist {
			return []int{haveSeen[another], index_one}
		}
		haveSeen[one] = index_one
	}
	return nil
}
