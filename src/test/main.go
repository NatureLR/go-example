package main

import (
	"flag"
	"fmt"

	"git.likeit.cn/go/audit"
)

func main() {
	ver := flag.Bool("version", false, "show version info")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}

	fmt.Println("aaaa")
	return
	audit.SetDebugging(true)
	getData()
}

func sortArr(arr []int, size int) []int {
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
