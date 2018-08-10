package main

import (
	"database/sql"
	"fmt"

	"git.likeit.cn/go/aux"
)

func doCreateDate(db *sql.DB) {
	tx, err := db.Begin()
	assert(err)
	defer func() {
		assert(tx.Commit())
	}()

	for i := 1; i < 500; i++ {
		sid := "102890"
		orders_id := aux.UUID(16)
		_, err := tx.Exec(`REPLACE INTO orders(shop_id,id,deleted,cust_name) VALUES (?,?,?,?)`, sid, orders_id, 1, "shopmove_test")
		assert(err)
	}
	fmt.Println("插入了500条数据")
}

func createData() {
	db, err := sql.Open("mysql", "root:renfeishengxian@tcp(192.168.70.228:3306)/shop?charset=utf8&sql_mode=ANSI_QUOTES")
	assert(err)
	defer db.Close()
	for i := 1; i < 500; i++ {
		doCreateDate(db)
	}
	fmt.Println()
}
