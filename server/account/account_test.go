package account

import (
	// "os"
	"encoding/json"
	"fmt"
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

	AddRecord := func(data rd.Item) error {
		r := Record{data}
		return r.Add()
	}
	DelRecord := func(data rd.Item) error {
		r := Record{data}
		return r.Del()
	}
	ChgRecord := func(data rd.Item) error {
		r := Record{data}
		return r.Chg()
	}

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
			v: rd.Item{ID: "16", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Account: [2]string{"CH1234", ""}, Amount: 60}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1078.5, Unit: 20.5}},
			d: []db.Debit{{ID: "0ef6832705b248cacd73b860cc95ee18693a5712", Name: "li4", Amount: 60}},
		},
		{ // chg b
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "ec6013a35aa2990b5f546e8862da8991066bd447", Type: rd.TypeB, Time: time.Date(2016, 12, 25, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"MO6677", ""}, Unit: 5, Amount: 10, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1028.5, Unit: 10.5}, {ID: "MO6677", Type: "OK", Input: 210.5, Unit: 49}},
			d: []db.Debit{{ID: "c061fdc19522830d3f1c468eead59e507f35e40a", Name: "zhang3", Amount: 10, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		{ // chg b
			ok: false, f: ChgRecord,
			v: rd.Item{ID: "ec6013a35aa2990b5f546e8862da8991066bd447", Type: rd.TypeB, Time: time.Date(2016, 12, 25, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"MO6677", ""}, Unit: 5, Amount: 10, Deadline: time.Date(2017, 1, 1, 1, 0, 0, 0, time.Local)}, n: 4,
		},
		{ // del b
			ok: true, f: DelRecord,
			v: rd.Item{ID: "95b65a60dab4b1872fcf11db9da04ac3e914895c", Type: rd.TypeB, Time: time.Date(2016, 12, 25, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"MO6677", ""}, Unit: 5, Amount: 10, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}, n: 3,
			a: []db.Account{{ID: "MO6677", Type: "OK", Input: 200.5, Unit: 44}},
			d: []db.Debit{{ID: "c061fdc19522830d3f1c468eead59e507f35e40a", Note: "Del"}},
		},
		{ // chg l
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "7bc80a2afbb7ce2962604380f1acd218eb02b38b", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Member: "li4", Account: [2]string{"CH1234", ""}, Amount: 70}, n: 3,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1018.5, Unit: 10.5}},
			d: []db.Debit{{ID: "0ef6832705b248cacd73b860cc95ee18693a5712", Name: "li4", Amount: 70}},
		},
		{ // del l
			ok: true, f: DelRecord,
			v: rd.Item{ID: "1f2e6aaf03a8f2c4804d40db7b5f3bc8b86c6547", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Member: "li4", Account: [2]string{"CH1234", ""}, Amount: 70}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: 1088.5, Unit: 10.5}},
			d: []db.Debit{{ID: "0ef6832705b248cacd73b860cc95ee18693a5712", Note: "Del"}},
		},
	} {
		if len(tc.v.ID) == 0 { // get & set
			for _, v := range buf {
				fmt.Println(v)
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
			v, err := tab.Sel(d)
			if err != nil {
				if !tc.ok || d.Note == "Del" {
					continue
				}
				fmt.Println(d.Note)
				t.Fatalf("%d,%d:%v", i, j, err)
			}
			if d.Note == "Del" {
				t.Fatalf("%d,%d:%v", i, j, "should be del")
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

	records, err := Records("2016-01-01 00:00:00", "2017-12-01 00:00:00")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatal(len(records))
	}
	// fmt.Println(records)
}

func TestRecordJson(t *testing.T) {
	for i, tc := range []struct {
		v rd.Item
	}{
		{rd.Item{ID: "id1234",
			Type:     rd.TypeR,
			Time:     time.Date(2016, 12, 21, 12, 3, 4, 0, time.Local),
			Class:    [2]string{"ABC", "456"},
			Account:  [2]string{"CH1234", "MH6677"},
			Unit:     10,
			Amount:   100.5,
			Proj:     "May",
			Deadline: time.Date(2016, 12, 21, 12, 3, 4, 0, time.Local),
			Note:     "hello",
			Member:   "JJ",
		}},
		{rd.Item{ID: "id1234",
			Type:    rd.TypeR,
			Time:    time.Date(2016, 12, 21, 12, 3, 4, 0, time.Local),
			Account: [2]string{"CH1234"},
			Amount:  100.5,
		}},
		{rd.Item{ID: "id1234",
			Type:    rd.TypeO,
			Time:    time.Date(2016, 12, 21, 12, 3, 4, 0, time.Local),
			Account: [2]string{"", "CH1234"},
			Class:   [2]string{"ABC"},
			Amount:  100.5,
		}},
		{rd.Item{ID: "id1234",
			Time:    time.Date(2016, 12, 21, 12, 3, 4, 0, time.Local),
			Account: [2]string{"", "CH1234"},
			Class:   [2]string{"", "ABC"},
		}},
	} {
		r := Record{tc.v}
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		// fmt.Println(string(b))
		var v Record
		if err := json.Unmarshal(b, &v); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(r, v) {
			t.Fatalf("%d: %+v,%+v", i, v, r)
		}
	}
}

func TestAccountIs(t *testing.T) {
	for i, tc := range []struct {
		a   Account
		isA bool
		isD bool
	}{
		{Account{}, false, false},
		{Account{Account: db.Account{ID: "1"}}, true, false},
		{Account{Debit: db.Debit{ID: "2"}}, false, true},
		{Account{db.Account{ID: "1"}, db.Debit{ID: "2"}}, true, true},
	} {
		if tc.a.isAccount() != tc.isA {
			t.Fatalf("%d: A %v", i, tc.a.isAccount())
		}
		if tc.a.isDebit() != tc.isD {
			t.Fatalf("%d: D %v", i, tc.a.isDebit())
		}
	}
}

func TestAccount(t *testing.T) {
	const (
		FMT = "2006-01-02 15:04:05"
	)

	min := time.Now()

	for i, tc := range []struct {
		ok bool
		v  Account
		n  int // record
		a  []db.Account
		d  []db.Debit
	}{
		{ // add a
			ok: true,
			v:  Account{Account: db.Account{ID: "CD9007", Type: "A1", Unit: 0.55, NUV: 0.34, Class: "C1", Input: 100.55, Deadline: time.Date(2017, 2, 1, 2, 23, 12, 22, time.Local)}}, n: 1,
			a: []db.Account{{ID: "CD9007", Type: "A1", Unit: 0.55, NUV: 0.34, Class: "C1", Input: 100.55, Deadline: time.Date(2017, 2, 1, 2, 23, 12, 22, time.Local)}},
		},
		{ // add a repeat
			ok: false,
			v:  Account{Account: db.Account{ID: "CD9007", Type: "A1", Unit: 0.55, NUV: 0.34, Class: "C1", Input: 100.55, Deadline: time.Date(2017, 2, 1, 2, 23, 12, 22, time.Local)}}, n: 1,
		},
		{ // chg a
			ok: true,
			v:  Account{Account: db.Account{ID: "CD9007", Type: "A2", Class: "C2", Input: 200}}, n: 2,
			a: []db.Account{{ID: "CD9007", Type: "A2", Class: "C2", Input: 200}},
		},
		{ // del a
			ok: true,
			v:  Account{Account: db.Account{ID: "CD9007"}}, n: 3,
			a: []db.Account{{ID: "CD9007", Type: "Del"}},
		},
		{ // del a repeat
			ok: false,
			v:  Account{Account: db.Account{ID: "CD9007"}}, n: 3,
		},
	} {
		err := tc.v.Update()
		if err != nil { // record func
			if !tc.ok {
				continue
			}
			t.Fatalf("%d:%v", i, err)
		}
		buf, err := rd.Get(min.Format(FMT), time.Now().Format(FMT)) // get record
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
				if !tc.ok || a.Type == "Del" {
					continue
				}
				t.Fatalf("%d,%d:%v", i, j, err)
			}
			if a.Type == "Del" {
				t.Fatalf("%d,%d:%v", i, j, "should be del")
			}
			v.Time = a.Time
			if v.Deadline.Unix() != a.Deadline.Unix() {
				if !tc.ok {
					continue
				}
				t.Fatalf("%d,%d:%v", i, j, v.Deadline)
			}
			v.Deadline = a.Deadline
			if !reflect.DeepEqual(v, a) {
				if !tc.ok {
					continue
				}
				t.Fatalf("%d,%d:%v", i, j, v)
			}
		}
		for j, d := range tc.d {
			var tab db.DebitTable
			if tc.v.Debit.Type == db.Borrow {
				tab = db.GetBorrow()
			} else if tc.v.Debit.Type == db.Lend {
				tab = db.GetLend()
			}
			v, err := tab.Sel(d)
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
	// os.RemoveAll("db")
}
