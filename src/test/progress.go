package main

import (
	"fmt"
)

type progress struct {
	Module  string `json:"module"`   //模块名字如x_orders
	ShopID  string `json:"shop_id"`  //店铺id如104343
	StatKey string `json:"stat_key"` //条目的进度id如order_id
	StatUpd string `json:"stat_upd"` //数据的更新进度
	Updated string `json:"updated"`  //中转库上的更新时间
}

func progressive() {
	var p progress
	localprogress(&p)
	fmt.Println("*************", p.Module)
}

func localprogress(pr *progress) {
	tp := progress{
		Module: "test",
	}
	pr = &tp
}
