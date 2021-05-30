package array

import (
	"reflect"
	"testing"
)

func TestTwoSum(t *testing.T) {
	nums := []int{2, 7, 11, 15}
	target := 9
	want := []int{0, 1}
	result := twoSum(nums, target)
	if !reflect.DeepEqual(result, want) {
		t.Error("want", want, "result:", result)
	}

}
