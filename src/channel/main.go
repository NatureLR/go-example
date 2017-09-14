package main

import (
	"fmt"
	"strconv"
)

func main() {
	taskChan := make(chan string, 3)
	doneChan := make(chan int, 1)
	for i := 0; i < 3; i++ {
		taskChan <- strconv.Itoa(i)
		fmt.Println("发送: ", i)
	}
	go func() {
		for i := 0; i < 3; i++ {
			task := <-taskChan
			fmt.Println("接受: ", task)
		}
		doneChan <- 1
	}()
	<-doneChan
}
