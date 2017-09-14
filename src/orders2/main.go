package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const BATCH = 1000

func main() {
	os.MkdirAll("data", 0755)
	ver := flag.Bool("version", false, "show version info")
	shop := flag.String("shops", "", "Shop ID to export")
	root := flag.String("root", "data", "Root directory of exported data")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	rx := regexp.MustCompile(`^[a-zA-Z0-9]{16}$`)
	if !rx.Match([]byte(*shop)) {
		fmt.Println("Invalid shop-id:", *shop)
		return
	}
	st, err := os.Stat(*root)
	if err != nil || !st.IsDir() { //判断输入的参数是否为空和是否合法
		fmt.Println("ERROR: -root not specified or is not a valid directory")
		return
	}
	db, err = sql.Open("mysql", "zxz@tcp(192.168.8.81:3306)/main?charset=utf8&sql_mode=ANSI_QUOTES")
	assert(err)
	expOrders(*root, *shop)
	expOrdItems(*root, *shop)
	expToCSV(*root, *shop)
	execzip()
}
