package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

func expOrders(output, sid string) {
	//导出订单
	t := time.Now()
	fmt.Printf("[%v] 导出订单\n", t)
	offset := 0
	for {
		sql := `SELECT "order".id,"order"."$shop","order".src,"order".type AS 
		订单类型,"order".number AS 订单号,"order".amount AS 订单金额,"order".venue 
		AS "房台/牌号",date_format(FROM_UNIXTIME("order".closed),'%Y-%m-%d') AS 日期,
		date_format(FROM_UNIXTIME("order".closed), '%T') AS 时间,"order".closed,order_delivery.
		"$order",order_delivery.name AS 名字,order_delivery.contact AS 电话,order_delivery.
		address AS 地址 FROM "order" LEFT JOIN order_delivery ON "order".id=order_delivery."$order" 
		WHERE "order"."$stat"='cls' AND "order"."$shop"=?`
		rs, err := db.Query(sql+` LIMIT ? OFFSET ?`, sid, BATCH, offset)
		assert(err)
		data := FetchRows(rs)
		fmt.Println(data)
		offset += len(data)
		fmt.Print(".")
		oids := []interface{}{}
		for _, d := range data {
			switch d["订单类型"] {
			case "0":
				d["订单类型"] = "订单"
			default:
				d["订单类型"] = "退单"
			}
			d["订单来源"] = orderType(d)
			var rt []string
			if json.Unmarshal([]byte(d["房台/牌号"]), &rt) == nil {
				d["房台/牌号"] = strings.Join(rt, ",")
			} else {
				d["房台/牌号"] = ""
			}
			oids = append(oids, d["id"])
		}
		sql = `SELECT "$order",gname FROM order_payment WHERE "$order" IN (?` +
			strings.Repeat(`,?`, len(oids)-1) + `)`
		rs, err = db.Query(sql, oids...)
		assert(err)
		pays := make(map[string][]string)
		for _, p := range FetchRows(rs) {
			oid := p["$order"]
			pays[oid] = append(pays[oid], p["gname"])
		}
		for _, d := range data {
			pay := pays[d["id"]]
			d["支付方式"] = strings.Join(pay, "/")
			dir := path.Join(output, d["$shop"], "订单")
			if os.MkdirAll(dir, 0755) != nil {
				continue
			}
			fn := path.Join(dir, d["closed"]+"_"+d["id"]+".json")
			delete(d, "src")
			delete(d, "id")
			delete(d, "$shop")
			delete(d, "closed")
			delete(d, "$order")
			f, err := os.Create(fn)
			assert(err)
			enc := json.NewEncoder(f)
			enc.SetIndent("", "    ")
			enc.Encode(d)
			assert(f.Close())
		}
		if len(data) < BATCH {
			fmt.Println()
			break
		}
	}
	fmt.Printf("elapsed: %f seconds\n", time.Since(t).Seconds())
}
