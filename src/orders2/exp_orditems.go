package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

func expOrdItems(output, sid string) {
	t := time.Now()
	fmt.Printf("[%v] 导出订单明细\n", t)
	offset := 0
	//导出订单明细
	for {
		sql := `SELECT order_item.id AS oiid,"order"."$shop" AS shop_id,
		    "order".type AS 订单类型,"order".number AS 订单号,venue AS 房台,src,code
			AS 商品码,name AS 品名,price AS 价格,count AS 数量,cat AS 分类,date_format(
			FROM_UNIXTIME("order".closed),'%Y-%m-%d') AS 日期,date_format(
			FROM_UNIXTIME("order".closed), '%T') AS 时间,"order".closed FROM
			order_item,"order" WHERE order_item."$order"="order".id AND
			"order"."$stat"='cls' AND order_item."$shop"=?`
		rs, err := db.Query(sql+` LIMIT ? OFFSET ?`, sid, BATCH, offset)
		assert(err)
		data := FetchRows(rs)
		offset += len(data)
		fmt.Print(".")
		for _, d := range data {
			switch d["订单类型"] {
			case "0":
				d["订单类型"] = "订单"
			default:
				d["订单类型"] = "退单"
			}
			d["订单来源"] = orderType(d)
			var rt []string
			if json.Unmarshal([]byte(d["房台"]), &rt) == nil {
				d["房台"] = strings.Join(rt, ",")
			} else {
				d["房台"] = ""
			}
			dir := path.Join(output, d["shop_id"], "订单明细")
			assert(os.MkdirAll(dir, 0755))
			fn := path.Join(dir, d["closed"]+"_"+d["oiid"]+".json")
			delete(d, "src")
			delete(d, "oiid")
			delete(d, "shop_id")
			delete(d, "closed")
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
