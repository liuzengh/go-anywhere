package binsearch

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// 给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的 中位数 。
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	if len(nums2) < len(nums1) {
		return findMedianSortedArrays(nums2, nums1)
	}
	low, high, n := 0, len(nums1), len(nums1)+len(nums2)
	left_num := (n + 1) / 2
	var split1, split2 int
	for low <= high {
		split1 = (low + high) / 2
		split2 = left_num - split1
		if split1 < len(nums1) && nums1[split1] < nums2[split2-1] {
			low = split1 + 1
		} else if 0 < split1 && nums2[split2] < nums1[split1-1] {
			high = split1 - 1
		} else {
			break
		}
	}
	var left, right int
	if split1 == 0 {
		left = nums2[split2-1]
	} else if split2 == 0 {
		left = nums1[split1-1]
	} else {
		left = max(nums1[split1-1], nums2[split2-1])
	}
	if n%2 == 1 {
		return float64(left)
	}
	if split1 == len(nums1) {
		right = nums2[split2]
	} else if split2 == len(nums2) {
		right = nums1[split1]
	} else {
		right = min(nums1[split1], nums2[split2])
	}
	return float64(left+right) / 2
}

/*
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    nums := append(nums1, nums2...)
    sort.Ints(nums)
    return float64(nums[len(nums)/2] + nums[(len(nums)-1)/2]) / 2
}
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    var nums []int
    for i, j := 0, 0; i < len(nums1) || j < len(nums2); {
        if i < len(nums1) && j < len(nums2) {
            if nums1[i] < nums2[j] {
                nums =  append(nums, nums1[i])
                i++
            } else {
                nums = append(nums, nums2[j])
                j++
            }
        } else if i < len(nums1) {
            nums = append(nums, nums1[i])
            i++
        } else {
            nums = append(nums, nums2[j])
            j++
        }
    }
    return float64(nums[len(nums)/2] + nums[(len(nums)-1)/2]) / 2

}
*/
