package main

import "fmt"

// 定义接口
type iPeople interface {
	get() string
	put(name string, sex bool)
}

type people struct {
	name string
	sex  bool
}

// 实现了get接口
func (t *people) get() string {
	if t.sex {
		return "" + t.name + " sex: " + "男"
	}
	return "" + t.name + " sex: " + "女"
}

// 实现了put接口
func (t *people) put(name string, sex bool) {
	t.name = name
	t.sex = sex
}

// 返回是一个接口
func NewPeople(name string, sex bool) iPeople {
	return &people{
		name: name,
		sex:  sex,
	}
}

// 参数是一个接口
func peopelInfo(i iPeople) {
	fmt.Println("info", i.get())
}

func main() {
	x := NewPeople("zxz", true)
	fmt.Println(x.get())
	x.put("zxz", false)
	fmt.Println(x.get())
	peopelInfo(x)
}
