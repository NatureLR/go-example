package main

import (
	//	"encoding/csv"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type opdata struct {
	total      string
	pay_amount string
}
type localdata struct {
	total      string
	pay_amount string
	variance   string
}

func readfile(filename string) map[string]localdata {
	filedata := make(map[string]localdata)
	f, err := os.Open(filename)
	assert(err)
	defer f.Close()
	buff := bufio.NewReader(f)
	for {
		line, err := buff.ReadString('\n') //一行一行的读取
		if err != nil || io.EOF == err {
			break
		}
		line = strings.TrimSpace(line)
		word := strings.Split(line, ",")
		filedata[word[0]] = localdata{total: word[1], pay_amount: word[2], variance: word[3]}
	}
	return filedata
}
func getoids(f string) []string {
	var oids []string
	for oid, _ := range readfile(f) {
		if oid == "id" {
			continue
		}
		oids = append(oids, oid)
	}
	return oids
}

func main() {
	ver := flag.Bool("version", false, "show version info")
	file := flag.String("f", "", "要比对的文件")
	tp := flag.String("t", "", "payment:比对payment,item比对Orderitem")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	if *file == "" {
		fmt.Println("没有文件 请查看-h")
		return
	}
	if *tp == "payment" {
		ComparePayment(*file)
	} else if *tp == "item" {
		CompareItem(*file)
	} else {
		fmt.Println("请-h查看帮助")
	}
}
