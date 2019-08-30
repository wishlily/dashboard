package database

import (
	"os"
	"testing"
)

func TestDBAccount(t *testing.T) {
	const URL = "test.db"
	defer os.Remove(URL)

	v1 := Account{
		Time:     12313,
		ID:       "cd1234",
		Type:     "CNY",
		Unit:     123.12,
		NUV:      0.8,
		Class:    "open",
		Input:    150,
		Deadline: 654321,
	}
	v2 := v1
	v2.ID = "ff2567"
	db, err := newAccountTable(URL)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	// add
	if err := db.add(v1); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := db.add(v1); err == nil {
		t.Fatalf("Add should be error")
	}
	if err := db.add(v2); err != nil {
		t.Fatalf("Add : %v", err)
	}
	// chg
	v1.Input = 200
	if err := db.chg(v1); err != nil {
		t.Fatalf("Change: %v", err)
	}
	v, err := db.sel(v1.ID)
	v1.Time = v.Time // time
	if err != nil || v != v1 {
		t.Fatalf("Get: %v: %v", err, v)
	}
	// get
	vv, err := db.get()
	if err != nil {
		t.Fatalf("Get %v", err)
	}
	if len(vv) != 2 {
		t.Fatalf("Get shoudle be 2 : %v", len(vv))
	}
	// del
	if err := db.del(v1); err != nil {
		t.Fatalf("Del: %v", err)
	}
	vv, err = db.get()
	if err != nil || len(vv) != 1 {
		t.Fatalf("Data: %v,%v", err, vv)
	}
}
