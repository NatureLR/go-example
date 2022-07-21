package main

import (
	"fmt"
	"reflect"
)

type maker interface {
	get() string
	post() string
}

type user struct {
	name string
}

func (u *user) get() string {
	return u.name + "get"
}

func (u *user) post() string {
	return u.name + "post"
}

type admin struct {
	rule string
	maker
}

func NewAdm(a maker) admin {
	t := reflect.TypeOf(a).String()
	return admin{
		maker: a,
		rule:  "ddd" + t,
	}
}

func (a *admin) get() string {
	return a.maker.get() + "admin get"
}

func (a *admin) post() string {
	return a.maker.post() + "admin post"
}

func makelove(i maker) {
	fmt.Println(i.post())
	fmt.Println(i.get())
}

func makelove2(i admin) {
	fmt.Println(i.rule)
	fmt.Println(i.post())
	fmt.Println(i.get())
}

func Reload() {
	a := &user{name: "zhangsan"}
	makelove(a)

	x := NewAdm(a)
	makelove2(x)
}
