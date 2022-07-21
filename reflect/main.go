package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Name struct {
	Sex   bool
	Age   int `json:"age" yaml:"age"`
	class string
}

func (n Name) Print() {
	fmt.Println(n.Sex)
}

func main() {

	name := Name{
		Sex:   true,
		Age:   2,
		class: "上班",
	}

	tr := reflect.TypeOf(name)
	fmt.Printf("r.Kind(): %s\n", tr.Kind().String())
	fmt.Println(tr.MethodByName("print"))

	tv, y := tr.FieldByName("Age")
	fmt.Println(y)
	fmt.Println(tv.Offset)
	fmt.Println(tv.Tag.Lookup("json"))
	fmt.Println(tv.Type)

	vr := reflect.ValueOf(name)
	vv := vr.FieldByName("Age")

	fmt.Println("xxxx")
	for i := 0; i < vr.NumField(); i++ {
		fieldInfo := vr.Type().Field(i) // a reflect.StructField

		name := strings.ToLower(fieldInfo.Name)

		fmt.Println(name, vr.Field(i))

	}
	fmt.Printf("vv.Int(): %v\n", vv.Int())
}
