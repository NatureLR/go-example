package main

import "fmt"

type Item struct {
	Name string
}

func main() {

	items := []Item{{Name: "zxzz"}, {Name: "dengmeng"}}

	var all []*Item

	for _, item := range items {
		fmt.Println(&item.Name)
		//item := item
		all = append(all, &item)
	}

	for _, i := range all {
		fmt.Println(*i)
	}
}
