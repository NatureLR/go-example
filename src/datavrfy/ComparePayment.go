package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func paymentquery(orderids []string) map[string]map[string]opdata {
	dbdata := make(map[string]map[string]opdata)
	dsns := []string{
		"shop-01-ro-01:KaIL0BkKGuPfLTiC@tcp(unit-sql-01.paadoo.net:3306)/shop",
		"shop-02-ro-01:Me1ftRCV23LVe7Vo@tcp(unit-sql-02.paadoo.net:3306)/shop",
		"shop-03-ro-01:QBSH5KOFsdsC6Cc8@tcp(unit-sql-03.paadoo.net:3306)/shop",
		"shop-04-ro-01:v2OwExSyuqGBZ1GJ@tcp(unit-sql-04.paadoo.net:3306)/shop",
		"shop-05-ro-01:Js6ZRlytdSdSV7N3@tcp(unit-sql-05.paadoo.net:3306)/shop",
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
			o.id in ('` + strings.Join(orderids, "','") + `')`
		re, err := db.Query(qry)
		assert(err)
		data := FetchRows(re)
		for _, r := range data {
			if dbdata[r["order_id"]] == nil {
				dbdata[r["order_id"]] = make(map[string]opdata)
			}
			dbdata[r["order_id"]][r["payment_id"]] = opdata{
				total:      r["total"],
				pay_amount: diffs("0", r["pay_amount"], 2),
			}
		}
	}
	return dbdata
}

func ComparePayment(file string) {
	ldb := readfile(file)
	dbd := paymentquery(getoids(file))
	fmt.Println("order_id,total,pay_amount,variance,payment_id,toal,pay_amount")
	for _, oid := range getoids(file) {
		if oid == "id" {
			continue
		}
		if len(dbd[oid]) > 1 { //一个order_id 多个payment
			for pid, v := range dbd[oid] {
				fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", oid, ldb[oid].total, ldb[oid].pay_amount, ldb[oid].variance, pid, v.total, v.pay_amount)
			}
		} else { //一个order_ID单个paymentID
			for pid, v := range dbd[oid] {
				lt, _ := strconv.ParseFloat(ldb[oid].total, 64)
				vt, _ := strconv.ParseFloat(v.total, 64)
				ly, _ := strconv.ParseFloat(ldb[oid].pay_amount, 64)
				vp, _ := strconv.ParseFloat(v.pay_amount, 64)
				if lt == vt && ly == vp {
					fmt.Printf("%s,%s,%s,%s,%s,%0.2f,%0.2f,√\n", oid, ldb[oid].total, ldb[oid].pay_amount, ldb[oid].variance, pid, vt, vp)
				} else {
					fmt.Printf("%s,%s,%s,%s,%s,%0.2f,%0.2f\n", oid, ldb[oid].total, ldb[oid].pay_amount, ldb[oid].variance, pid, vt, vp)
				}
			}
		}
	}
}
