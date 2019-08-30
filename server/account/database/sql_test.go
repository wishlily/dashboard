package database

import (
	"fmt"
	"os"
	"testing"
)

func TestTable(t *testing.T) {
	const URL = "test.db"
	defer os.Remove(URL)

	type DB struct {
		/* TAG db MUST SET */
		OutradeNo string  `db:"out_trade_no" type:"VARCHAR(32)"`
		Valid     bool    `db:"flag" type:"TINYINT(1)"`
		Value     float64 `type:"DECIMAL(32,3)"`
		TradeNo   string  `db:"-"`
		Time      int64   `db:"time" type:"BIGINT"`
	}
	v1 := DB{"654321", true, 233.23, "hello", 789}
	v2 := DB{"123", false, -12.4, "world", 55}
	db, err := newTable(URL, "test1", v1)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	// insert
	if err := db.insert(v1); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if err := db.insert(v2); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	// update
	v1.Value = 123.123
	v1.Time = 987
	if err := db.update("out_trade_no", v1); err != nil {
		t.Fatalf("Update: %v", err)
	}
	rows, err := db.sel("time", v1)
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	v1.TradeNo = "" // not store
	for rows.Next() {
		var m DB
		err = rows.StructScan(&m)
		if err != nil {
			t.Fatalf("StructScan : %v", err)
		}
		if m != v1 {
			t.Fatalf("Equal : %+v,%+v", m, v1)
		}
	}
	// get
	rows, err = db.get()
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	num := 0
	for rows.Next() {
		var m DB
		err = rows.StructScan(&m)
		if err != nil {
			t.Fatalf("StructScan : %v", err)
		}
		num++
		fmt.Println(m)
	}
	if num != 2 {
		t.Fatalf("Get should be 2 : %v", num)
	}
	// delete
	if err := db.delete("out_trade_no", v1); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := db.update("out_trade_no", v1); err == nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := db.drop(); err != nil {
		t.Fatalf("Drop: %v", err)
	}
}

func TestTableErr(t *testing.T) {
	const URL = "test.db"
	defer os.Remove(URL)

	if _, err := newTable(URL, "", "hello"); err == nil {
		t.Fatalf("new need error")
	}
}
