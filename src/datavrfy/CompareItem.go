package main

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func itemquery(orderids []string) map[string]map[string]opdata {
	//fmt.Println("查询数据库中稍等....")
	dbdata := make(map[string]map[string]opdata)
	dsns := []string{
		"shop-01-ro-01:KaIL0BkKGuPfLTiC@tcp(unit-sql-01.paadoo.net:3306)/shop",
		"shop-02-ro-01:Me1ftRCV23LVe7Vo@tcp(unit-sql-02.paadoo.net:3306)/shop",
		"shop-03-ro-01:QBSH5KOFsdsC6Cc8@tcp(unit-sql-03.paadoo.net:3306)/shop",
		"shop-04-ro-01:v2OwExSyuqGBZ1GJ@tcp(unit-sql-04.paadoo.net:3306)/shop",
		"shop-05-ro-01:Js6ZRlytdSdSV7N3@tcp(unit-sql-05.paadoo.net:3306)/shop",
	}
	for _, dsn := range dsns {
		//fmt.Print(".")
		db, err := sql.Open("mysql", dsn)
		assert(err)
		qry := `SELECT DISTINCT
					o.id AS order_id,
					oi.id AS item_id,
					o.total,
					oi.final_price,
					oi.counts
				FROM
					orders o
					INNER JOIN one_order oo ON o.id = oo.order_id
					INNER JOIN order_item oi ON oo.id = oi.one_order_id
				WHERE
					o.id in ('` + strings.Join(orderids, "','") + `')`
		re, err := db.Query(qry)
		assert(err)
		data := FetchRows(re)
		for _, r := range data {
			if dbdata[r["order_id"]] == nil {
				dbdata[r["order_id"]] = make(map[string]opdata)
			}
			dbdata[r["order_id"]][r["item_id"]] = opdata{
				total:      r["total"],
				pay_amount: prods(r["final_price"], r["counts"], 2),
			}
		}
		//fmt.Println()
	}
	return dbdata
}
func CompareItem(file string) {
	ldb := readfile(file)
	dbd := itemquery(getoids(file))
	fmt.Println("order_id,total,final_amount,variance,数据库的toal,数据库的final_amount")
	for _, oid := range getoids(file) {
		if oid == "id" {
			continue
		}
		var t float64 = 0
		var f float64 = 0
		for _, v := range dbd[oid] {
			t, _ = strconv.ParseFloat(v.total, 64)
			vp, _ := strconv.ParseFloat(v.pay_amount, 64)
			f = f + vp
		}
		lt, _ := strconv.ParseFloat(ldb[oid].total, 64)
		ly, _ := strconv.ParseFloat(ldb[oid].pay_amount, 64)
		if math.Abs(lt-t) < 0.005 && math.Abs(ly-f) < 0.005 {
			fmt.Printf("%s,%0.2f,%0.2f,%s,%0.2f,%0.2f,√\n", oid, lt, ly, ldb[oid].variance, t, f)
		} else {
			for iid, v := range dbd[oid] {
				fmt.Printf("%s,%0.2f,%0.2f,%s,%s,%s,%s\n", oid, lt, ly, ldb[oid].variance, v.total, v.pay_amount, iid)
			}
		}
	}
}
