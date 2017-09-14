package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type (
	Row  map[string]string
	Rows []Row
)

func RangeRows(rows *sql.Rows, proc func()) {
	defer func() {
		if e := recover(); e != nil {
			rows.Close()
			panic(e)
		}
	}()
	for rows.Next() {
		proc()
	}
	assert(rows.Err())
}

func FetchRows(rows *sql.Rows) Rows {
	defer func() {
		if e := recover(); e != nil {
			rows.Close()
			panic(e)
		}
	}()
	cols, err := rows.Columns()
	assert(err)
	raw := make([][]byte, len(cols))
	ptr := make([]interface{}, len(cols))
	for i, _ := range raw {
		ptr[i] = &raw[i]
	}
	var recs Rows
	for rows.Next() {
		assert(rows.Scan(ptr...))
		rec := make(Row)
		for i, r := range raw {
			if r == nil {
				rec[cols[i]] = ""
			} else {
				rec[cols[i]] = string(r)
			}
		}
		recs = append(recs, rec)
	}
	assert(rows.Err())
	return recs
}
