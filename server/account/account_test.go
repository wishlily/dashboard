package account

import (
	"os"
	"reflect"
	"testing"
	"time"

	db "github.com/wishlily/dashboard/server/account/database"
	rd "github.com/wishlily/dashboard/server/account/record"
)

func TestInit(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
}

func TestRecordI(t *testing.T) {
	const (
		FMT = "2006-01-02 15:04:05"
	)
	ts := time.Now()

	a1 := db.Account{ID: "CH1234", Type: "FREE", Input: 1100.5}
	m1 := rd.Item{Type: rd.TypeI, Time: ts, Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 100.5}
	if err := db.GetAccount().Add(a1); err != nil {
		t.Fatal(err)
	}

	// Add
	if err := AddRecord(m1); err != nil {
		t.Fatal(err)
	}
	a1.Input += m1.Amount
	a1.Unit += m1.Unit
	{ // check record
		v, err := rd.Get(ts.Format(FMT), ts.Format(FMT))
		if err != nil {
			t.Fatal(err)
		}
		if len(v) != 1 {
			t.Fatal(err)
		}
		m1 = v[0]
	}
	{ // check account
		v, err := db.GetAccount().Sel(m1.Account[rd.AccountM])
		if err != nil {
			t.Fatal(err)
		}
		a1.Time = v.Time
		if !reflect.DeepEqual(a1, v) {
			t.Fatalf("%v", v)
		}
	}
}

func TestRemove(t *testing.T) {
	os.RemoveAll("db")
}
