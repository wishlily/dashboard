package database

import (
	"os"
	"testing"
	"time"
)

func TestDBAccount(t *testing.T) {
	const URL = "test.db"
	defer os.Remove(URL)

	v1 := Account{
		Time:  time.Date(2005, 1, 1, 12, 0, 0, 0, time.Local),
		ID:    "cd1234",
		Type:  "CNY",
		Unit:  123.12,
		NUV:   0.8,
		Class: "open",
		Input: 150,
	}
	v1.Deadline, _ = time.Parse("2006-01-02 15:04:05", "2006-01-01 12:00:00")
	v2 := v1
	v2.ID = "ff2567"
	db, err := newAccountTable(URL, "test")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	// add
	if err := db.Add(v1); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := db.Add(v1); err == nil {
		t.Fatalf("Add should be error")
	}
	if err := db.Add(v2); err != nil {
		t.Fatalf("Add : %v", err)
	}
	// chg
	v1.Input = 200
	if err := db.Chg(v1); err != nil {
		t.Fatalf("Change: %v", err)
	}
	v, err := db.Sel(v1.ID)
	v1.Time = v.Time // time
	if err != nil || v != v1 {
		t.Fatalf("Get: %v: %v,%v", err, v, v1)
	}
	// get
	vv, err := db.Get()
	if err != nil {
		t.Fatalf("Get %v", err)
	}
	if len(vv) != 2 {
		t.Fatalf("Get shoudle be 2 : %v", len(vv))
	}
	// del
	if err := db.Del(v1); err != nil {
		t.Fatalf("Del: %v", err)
	}
	vv, err = db.Get()
	if err != nil || len(vv) != 1 {
		t.Fatalf("Data: %v,%v", err, vv)
	}
}
