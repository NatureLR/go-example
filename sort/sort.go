package sort

//冒泡排序
func main(arr []int, size int) []int {
	for i := 0; i < size; i++ {
		for j := 0; j < (size - 1 - i); j++ {
			if arr[j] > arr[j+1] {
				tmp := arr[j+1]
				arr[j+1] = arr[j]
				arr[j] = tmp
			}
		}
	}
	return arr
}
