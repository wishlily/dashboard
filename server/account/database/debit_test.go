package database

import (
	"os"
	"testing"
)

func TestDBDebit(t *testing.T) {
	const URL = "test.db"
	defer os.Remove(URL)

	v1 := Debit{
		Name:     "zhang3",
		Amount:   200.02,
		Note:     "replace",
		Deadline: 123456,
	}
	v2 := Debit{
		Name:   "li4",
		Amount: 98,
	}
	db, err := newDebitTable(URL, "test3")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	// add
	if err := db.Add(v1); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := db.Add(v2); err != nil {
		t.Fatalf("Add : %v", err)
	}
	// get
	vv, err := db.Get()
	if err != nil {
		t.Fatalf("Get %v", err)
	}
	if len(vv) != 2 {
		t.Fatalf("Get shoudle be 2 : %v", len(vv))
	}
	v1 = vv[0]
	v2 = vv[1]
	// chg
	v1.Amount = 200
	if err := db.Chg(v1); err != nil {
		t.Fatalf("Change: %v", err)
	}
	if _, err := db.Sel(v1.ID); err != nil {
		t.Fatalf("Get: %v", err)
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
