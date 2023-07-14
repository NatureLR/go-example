package main

import (
	"bytes"
	"strings"
	"testing"
	"unsafe"
)

var L = 1024 * 1024
var str = strings.Repeat("a", L)
var s = bytes.Repeat([]byte{'a'}, L)
var str2 string
var s2 []byte

func BenchmarkString2Slice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bt := []byte(str)
		if len(bt) != L {
			b.Fatal()
		}
	}
}
func BenchmarkString2SliceReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bt := *(*[]byte)(unsafe.Pointer(&str))
		if len(bt) != L {
			b.Fatal()
		}
	}
}
func BenchmarkString2SliceUnsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bt := unsafe.Slice(unsafe.StringData(str), len(str))
		if len(bt) != L {
			b.Fatal()
		}
	}
}
func BenchmarkSlice2String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ss := string(s)
		if len(ss) != L {
			b.Fatal()
		}
	}
}
func BenchmarkSlice2StringReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ss := *(*string)(unsafe.Pointer(&s))
		if len(ss) != L {
			b.Fatal()
		}
	}
}
func BenchmarkSlice2StringUnsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ss := unsafe.String(unsafe.SliceData(s), len(str))
		if len(ss) != L {
			b.Fatal()
		}
	}
}
