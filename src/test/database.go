package main

import (
	"database/sql"
	"fmt"

	"git.likeit.cn/go/audit"
	"git.likeit.cn/go/aux"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	d, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8&sql_mode=ANSI_QUOTES&allowAllFiles=true")
	assert(err)
	db = d
}

//从mysql数据库中读取数据
func getData() {
	rows, err := db.Query("SELECT * FROM $policy")
	audit.Assert(err)

	/*for _, r := range aux.FetchRows(rows) {
		fmt.Println(r["id"])
	}

	aux.RangeRows(rows, func() {
		var f string
		rows.Scan(&f)
	})

	var count int
	test := func(r map[string]interface{}) bool {

		if count == 10 {
			return false
		}
		return true
	}
	aux.IterRows(rows, test)*/

	IterRowsLimit(rows, 10, func(r map[string]interface{}) {
		fmt.Println(r["id"])
	})
}

type proc func(map[string]interface{})

func IterRowsLimit(rows *sql.Rows, batch int, p proc) {
	var count int
	obj := func(r map[string]interface{}) bool {
		p(r)
		if count == batch {
			return false
		}
		count++
		return true
	}
	aux.IterRows(rows, obj)
}

//将本地文件导入到mysql中
func lif() {
	format := `(order_scene,is_takeout,member_name,amount,consumer_addr,position_type,pay2_amount,temp_high,shop_name,close_time,ext_order_num,pay1_method,parent_id,source,shop_tables,order_num,member_number,member_mobile,consumer_contact,shop_area1,shop_area2,order_date,member_gender,weekday,related,cust_name,area_name,weather,root_id,shop_lat,total,heads,lat,position_name,sessions,booking_time,service_time_slot,lon,pay1_opname,pay3_method,member_create_time,shop_seats,serial_num,remark,pay1_amount,holiday,temp_low,member_expiry,shop_area3,shop_lon,duration,mobile,pay3_amount,pay3_opname,$deleted,pos_owner_name,created,order_time,member_on_account,pay2_opname,table_type,order_type,device,member_id,member_account_paid,commission,pay2_method,id,updated,shop_id,member_type_name,consumer_name)`
	filePath := "/home/zxz/git/test/src/test/test.csv"
	mysql.RegisterLocalFile(filePath)
	stmt := `LOAD DATA LOCAL INFILE '` + filePath + `' INTO TABLE bdp_orders
	CHARACTER SET utf8  fields terminated by ',' enclosed by '"' ` + format
	fmt.Println(stmt)
	_, err := db.Exec(stmt)
	assert(err)
	return
	//DSN加上allowAllFiles=true
	//fields关键字指定了文件记段的分割格式，如果用到这个关键字，MySQL剖析器希望看到至少有下面的一个选项：
	//terminated by分隔符：意思是以什么字符作为分隔符
	//enclosed by字段括起字符
	//escaped by转义字符
	//terminated by描述字段的分隔符，默认情况下是tab字符（\t）
	//enclosed by描述的是字段的括起字符。
	//escaped by描述的转义字符。默认的是反斜杠（backslash：\ ）
}
