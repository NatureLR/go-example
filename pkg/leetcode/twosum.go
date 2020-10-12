package leetcode

func twoSum(nums []int, target int) []int {
	for k, v := range nums {
		for k1, v1 := range nums {
			v1 = v + v1
			if k == k1 {
				continue
			}
			if v1 == target {
				return []int{k, k1}
			}
		}
	}
	return nil
}
