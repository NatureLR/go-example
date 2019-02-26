package main

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//初始化数据库
func initDB() {
	d, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8&sql_mode=ANSI_QUOTES&allowAllFiles=true")
	assert(err)
	db = d
}

//创建数据库表
func createTable() {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS test (
		id int(11) NOT NULL AUTO_INCREMENT,
		created datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
		updated datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
		PRIMARY KEY (shop_id, module)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='进度控制'
	`)
	assert(err)
}

//从mysql数据库中读取数据
func getData() {
	//适合有不知道那些字段或者字段很多情况
	rows, err := db.Query("SELECT * FROM test")
	assert(err)
	for _, r := range FetchRows(rows) {
		fmt.Println(r["id"])
	}

	//适合有明确字段拿出来的情况
	rows, err = db.Query("SELECT id,created FROM test")
	RangeRows(rows, func() {
		var id string
		var create sql.NullString
		rows.Scan(&id, &create)
		if create.Valid {
			fmt.Println(id, create.String)
		}
	})

	//适合只有一条数据的情况
	var id sql.NullString //此对象为空时不会报错
	assert(db.QueryRow(`SELECT id FROM test`).Scan(&id))
}

//将本地文件导入到mysql中
func lif() {
	format := `(order_scene,is_takeout,member_name,amount,consumer_addr,position_type,pay2_amount,temp_high,shop_name,close_time,ext_order_num,pay1_method,parent_id,source,shop_tables,order_num,member_number,member_mobile,consumer_contact,shop_area1,shop_area2,order_date,member_gender,weekday,related,cust_name,area_name,weather,root_id,shop_lat,total,heads,lat,position_name,sessions,booking_time,service_time_slot,lon,pay1_opname,pay3_method,member_create_time,shop_seats,serial_num,remark,pay1_amount,holiday,temp_low,member_expiry,shop_area3,shop_lon,duration,mobile,pay3_amount,pay3_opname,$deleted,pos_owner_name,created,order_time,member_on_account,pay2_opname,table_type,order_type,device,member_id,member_account_paid,commission,pay2_method,id,updated,shop_id,member_type_name,consumer_name)`
	filePath := "/home/zxz/git/test/src/test/test.csv"
	mysql.RegisterLocalFile(filePath)
	stmt := `LOAD DATA LOCAL INFILE '` + filePath + `' INTO TABLE bdp_orders
	CHARACTER SET utf8 fields terminated by ',' enclosed by '"' ` + format

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
