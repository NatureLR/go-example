package main

import (
	"database/sql"
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
func dbquery(orderid []string) map[string]map[string]opdata {
	dbdata := make(map[string]map[string]opdata)
	dsns := []string{
		"shop-01-ro-01:KaIL0BkKGuPfLTiC@tcp(unit-sql-01.paadoo.net:3306)/shop",
		"shop-02-ro-01:Me1ftRCV23LVe7Vo@tcp(unit-sql-02.paadoo.net:3306)/shop",
		"shop-03-ro-01:QBSH5KOFsdsC6Cc8@tcp(unit-sql-03.paadoo.net:3306)/shop",
		"shop-04-ro-01:v2OwExSyuqGBZ1GJ@tcp(unit-sql-04.paadoo.net:3306)/shop",
	}
	for _, dsn := range dsns {
		db, err := sql.Open("mysql", dsn)
		assert(err)
		qry := `SELECT DISTINCT
				o.id AS order_id,
				p.id AS payment_id,
				o.amount AS total,
				p.paid AS pay_amount 
			FROM
				orders o
				INNER JOIN one_order oo ON o.id = oo.order_id
				INNER JOIN order_payment op ON op.one_order_id = oo.id
				INNER JOIN payment p ON p.id = op.payment_id 
			WHERE
			o.id in ('` + strings.Join(orderid, "','") + `')`
		re, err := db.Query(qry)
		assert(err)
		data := FetchRows(re)
		for _, r := range data {
			if dbdata[r["order_id"]] == nil {
				dbdata[r["order_id"]] = make(map[string]opdata)
			}
			dbdata[r["order_id"]][r["payment_id"]] = opdata{
				total:      r["total"],
				pay_amount: r["pay_amount"],
			}
		}
	}
	return dbdata
}

/*func ComparePayment(file string) {
	ldb := readfile(file)
	dbd := dbquery(getoids(file))
	fmt.Println("order_id,total,pay_amount,variance,payment_id,toal,pay_amount")
	for _, oid := range getoids(file) {
		if oid == "id" {
			continue
		}
		if ldb[oid].total == dbd[oid].total && ldb[oid].pay_amount == dbd[oid].pay_amount {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s,√\n", oid, ldb[oid].total, ldb[oid].pay_amount, ldb[oid].variance, dbd[oid].payment_id, dbd[oid].total, dbd[oid].pay_amount)
		} else {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", oid, ldb[oid].total, ldb[oid].pay_amount, ldb[oid].variance, dbd[oid].payment_id, dbd[oid].total, dbd[oid].pay_amount)
		}
	}
}*/

func main() {
	ver := flag.Bool("version", false, "show version info")
	file := flag.String("f", "test.csv", "要比对的文件")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	if *file == "" {
		fmt.Println("没有文件 请查看-h")
		return
	}

	ldb := readfile(*file)
	dbd := dbquery(getoids(*file))
	fmt.Println("order_id,total,pay_amount,variance,payment_id,toal,pay_amount")
	for _, oid := range getoids(*file) {
		if oid == "id" {
			continue
		}
		if len(dbd[oid]) > 1 { //一个order_id 多个payment
			for pid, v := range dbd[oid] {
				fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", oid, ldb[oid].total, ldb[oid].pay_amount, ldb[oid].variance, pid, v.total, v.pay_amount)
			}
		} else { //一个order_ID单个paymentID
			for pid, v := range dbd[oid] {
				if ldb[oid].total == v.total && ldb[oid].pay_amount == v.pay_amount {
					fmt.Printf("%s,%s,%s,%s,%s,%s,%s,√\n", oid, ldb[oid].total, ldb[oid].pay_amount, ldb[oid].variance, pid, v.total, v.pay_amount)
				} else {
					fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", oid, ldb[oid].total, ldb[oid].pay_amount, ldb[oid].variance, pid, v.total, v.pay_amount)
				}
				fmt.Println(oid, pid, v.total, v.pay_amount)
			}
		}
	}
}
