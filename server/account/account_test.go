package account

import (
	// "os"
	"fmt"
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
	a1 := db.Account{ID: "CH1234", Type: "FREE", Input: 1100.5}
	if err := db.GetAccount().Add(a1); err != nil {
		t.Fatal(err)
	}
	a2 := db.Account{ID: "MO6677", Type: "OK", Input: 100, Unit: 34}
	if err := db.GetAccount().Add(a2); err != nil {
		t.Fatal(err)
	}
}

func TestRecord(t *testing.T) {
	const (
		FMT = "2006-01-02 15:04:05"
	)

	min := time.Date(2016, 12, 21, 12, 3, 4, 11, time.Local)
	max := time.Date(2017, 2, 1, 2, 23, 12, 22, time.Local)
	var buf []rd.Item

	for i, tc := range []struct {
		ok bool
		f  func(item rd.Item) error
		v  rd.Item
		n  int
		a  []db.Account
		d  []db.Debit
	}{
		{ // add i
			ok: true, f: AddRecord,
			v: rd.Item{ID: "0", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 100.5}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1201, Unit: 10}},
		},
		{ // add o
			ok: true, f: AddRecord,
			v: rd.Item{ID: "1", Type: rd.TypeO, Time: min, Account: [2]string{"CH1234", ""}, Unit: 5.5, Amount: 12.23}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1188.77, Unit: 4.5}},
		},
		{ // add r
			ok: true, f: AddRecord,
			v: rd.Item{ID: "2", Type: rd.TypeR, Time: time.Date(2016, 12, 27, 10, 3, 4, 11, time.Local), Account: [2]string{"CH1234", "MO6677"}, Unit: 10, Amount: 100}, n: 3,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1088.77, Unit: 0}, {ID: "MO6677", Type: "OK", Input: 200, Unit: 44}},
		},
		{ // add x & error
			ok: false, f: AddRecord,
			v: rd.Item{ID: "3", Type: rd.TypeX}, n: 0,
			a: []db.Account{},
		},
		{ // del i
			ok: true, f: DelRecord,
			v: rd.Item{ID: "6249131937fc571f393b9f11c0e4168e5f47fde9", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 100.5}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 988.27, Unit: 0}},
		},
		{ // del o
			ok: true, f: DelRecord,
			v: rd.Item{ID: "33154f59e4081a230cb5ce3278a55caf7cd2ef6b", Type: rd.TypeO, Time: min, Account: [2]string{"CH1234", ""}, Unit: 5.5, Amount: 12.23}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1000.5, Unit: 5.5}},
		},
		{ // del r
			ok: true, f: DelRecord,
			v: rd.Item{ID: "ef993d2c4e4e327ae7469a91e18ea50561acde46", Type: rd.TypeR, Time: time.Date(2016, 12, 27, 10, 3, 4, 11, time.Local), Account: [2]string{"CH1234", "MO6677"}, Unit: 10, Amount: 100}, n: 0,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1100.5, Unit: 15.5}, {ID: "MO6677", Type: "OK", Input: 100, Unit: 34}},
		},
		{ // del x & error
			ok: false, f: DelRecord,
			v: rd.Item{ID: "7", Type: rd.TypeX}, n: 0,
			a: []db.Account{},
		},
		{ // add i
			ok: true, f: AddRecord,
			v: rd.Item{ID: "8", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 100.5, Note: "hello", Member: "jj", Proj: "xx"}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1201, Unit: 25.5}},
		},
		{ // chg i
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "cc35e99c2b4387b9e5db52782266475d84444947", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 100.5, Note: "hello", Member: "jj", Proj: "xx"}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1201, Unit: 25.5}},
		},
		{ // chg i
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "cc35e99c2b4387b9e5db52782266475d84444947", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"MO6677", ""}, Unit: 10, Amount: 100.5, Note: "hello", Member: "jj", Proj: "xx"}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1100.5, Unit: 15.5}, {ID: "MO6677", Type: "OK", Input: 200.5, Unit: 44}},
		},
		{ // add o
			ok: true, f: AddRecord,
			v: rd.Item{ID: "11", Type: rd.TypeO, Time: min, Account: [2]string{"CH1234", ""}, Unit: 5.5, Amount: 12.23, Member: "cc"}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1088.27, Unit: 10}},
		},
		{ // chg o
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "e54141dcbe8218b324798f42be4f9039103e512b", Type: rd.TypeO, Time: time.Date(2017, 1, 5, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 5, Amount: 12, Member: "bb"}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1088.5, Unit: 10.5}},
		},
		{ // add b
			ok: true, f: AddRecord,
			v: rd.Item{ID: "13", Type: rd.TypeB, Time: time.Date(2016, 12, 24, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 50, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}, n: 3,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1138.5, Unit: 20.5}},
			d: []db.Debit{{ID: "c061fdc19522830d3f1c468eead59e507f35e40a", Name: "zhang3", Amount: 50, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		{ // add b
			ok: false, f: AddRecord,
			v: rd.Item{ID: "14", Type: rd.TypeB, Time: time.Date(2016, 12, 24, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 50, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}, n: 3,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1138.5, Unit: 20.5}},
			d: []db.Debit{{ID: "c061fdc19522830d3f1c468eead59e507f35e40a", Name: "zhang3", Amount: 50, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		{ // add l
			ok: true, f: AddRecord,
			v: rd.Item{ID: "15", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Member: "li4", Account: [2]string{"CH1234", ""}, Amount: 60}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1078.5, Unit: 20.5}},
			d: []db.Debit{{ID: "0ef6832705b248cacd73b860cc95ee18693a5712", Name: "li4", Amount: 60}},
		},
		{ // add l
			ok: false, f: AddRecord,
			v: rd.Item{ID: "15", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Account: [2]string{"CH1234", ""}, Amount: 60}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1078.5, Unit: 20.5}},
			d: []db.Debit{{ID: "0ef6832705b248cacd73b860cc95ee18693a5712", Name: "li4", Amount: 60}},
		},
	} {
		if len(tc.v.ID) == 0 { // get & set
			for _, v := range buf {
				fmt.Println(v.ID)
			}
		}
		err := tc.f(tc.v)
		if err != nil { // record func
			if !tc.ok {
				continue
			}
			t.Fatalf("%d:%v", i, err)
		}
		buf, err = rd.Get(min.Format(FMT), max.Format(FMT)) // get record
		if err != nil {
			if !tc.ok {
				continue
			}
			t.Fatalf("%d:%v", i, err)
		}
		if len(buf) != tc.n {
			if !tc.ok {
				continue
			}
			t.Fatalf("%d:%d", i, len(buf))
		}
		for j, a := range tc.a {
			v, err := db.GetAccount().Sel(a.ID) // get account
			if err != nil {
				if !tc.ok {
					continue
				}
				t.Fatalf("%d,%d:%v", i, j, err)
			}
			v.Time = a.Time
			if !reflect.DeepEqual(v, a) {
				if !tc.ok {
					continue
				}
				t.Fatalf("%d,%d:%v", i, j, v)
			}
		}
		for j, d := range tc.d {
			var tab db.DebitTable
			if tc.v.Type == rd.TypeB {
				tab = db.GetBorrow()
			} else if tc.v.Type == rd.TypeL {
				tab = db.GetLend()
			}
			v, err := tab.Sel(d.ID)
			if err != nil {
				if !tc.ok {
					continue
				}
				t.Fatalf("%d,%d:%v", i, j, err)
			}
			v.Time = d.Time
			if v.Deadline.Unix() != d.Deadline.Unix() {
				if !tc.ok {
					continue
				}
				t.Fatalf("%d,%d:%v", i, j, v.Deadline)
			}
			v.Deadline = d.Deadline
			if !reflect.DeepEqual(v, d) {
				if !tc.ok {
					continue
				}
				t.Fatalf("%d,%d:%v", i, j, v)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	os.RemoveAll("db")
}
