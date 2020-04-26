package account

import (
	// "os"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	db "github.com/wishlily/dashboard/server/account/database"
	rd "github.com/wishlily/dashboard/server/account/record"
)

var (
	// CH1234, MO6677
	testAcct = []float64{1100.5, 100}
	testUnit = []float64{1, 34}
	testNuv  = []float64{1, 0}
)

func TestInit(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	a1 := db.Account{ID: "CH1234", Type: "FREE", Input: testAcct[0], Unit: testUnit[0], NUV: testNuv[0]}
	if err := db.GetAccount().Add(a1); err != nil {
		t.Fatal(err)
	}
	a2 := db.Account{ID: "MO6677", Type: "OK", Input: testAcct[1], Unit: testUnit[1], NUV: testNuv[1]}
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

	AddAcct := func(i int, v float64) float64 {
		testAcct[i] += v
		return testAcct[i]
	}

	SubAcct := func(i int, v float64) float64 {
		testAcct[i] -= v
		return testAcct[i]
	}

	AddUnit := func(i int, v float64) float64 {
		testUnit[i] += v
		return testUnit[i]
	}

	SubUnit := func(i int, v float64) float64 {
		testUnit[i] -= v
		if testUnit[i] < 0 {
			testUnit[i] = 0
		}
		return testUnit[i]
	}

	AddNuv := func(i int, v float64) float64 {
		testNuv[i] += v
		return testNuv[i]
	}

	SubNuv := func(i int, v float64) float64 {
		testNuv[i] -= v
		if testNuv[i] < 0 {
			testNuv[i] = 0
		}
		return testNuv[i]
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
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 100.5), Unit: AddUnit(0, 10), NUV: AddNuv(0, 0)}},
		},
		{ // add o
			ok: true, f: AddRecord,
			v: rd.Item{ID: "1", Type: rd.TypeO, Time: min, Account: [2]string{"CH1234", ""}, Unit: 5.5, Amount: 12.23}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 12.23), Unit: SubUnit(0, 5.5), NUV: SubNuv(0, 0)}},
		},
		{ // add r
			ok: true, f: AddRecord,
			v: rd.Item{ID: "2", Type: rd.TypeR, Time: time.Date(2016, 12, 27, 10, 3, 4, 11, time.Local), Account: [2]string{"CH1234", "MO6677"}, Unit: 1.1, Amount: 100}, n: 3,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 100), Unit: SubUnit(0, 1.1), NUV: SubNuv(0, 0)}, {ID: "MO6677", Type: "OK", Input: AddAcct(1, 100), Unit: AddUnit(1, 0), NUV: AddNuv(1, 0)}},
		},
		{ // add x & error
			ok: false, f: AddRecord,
			v: rd.Item{ID: "3", Type: rd.TypeX}, n: 0,
			a: []db.Account{},
		},
		{ // del i
			ok: true, f: DelRecord,
			v: rd.Item{ID: "6249131937fc571f393b9f11c0e4168e5f47fde9", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 100.5}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 100.5), Unit: SubUnit(0, 10), NUV: SubNuv(0, 0)}},
		},
		{ // del o
			ok: true, f: DelRecord,
			v: rd.Item{ID: "33154f59e4081a230cb5ce3278a55caf7cd2ef6b", Type: rd.TypeO, Time: min, Account: [2]string{"CH1234", ""}, Unit: 5.5, Amount: 12.23}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 12.23), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
		},
		{ // del r
			ok: true, f: DelRecord,
			v: rd.Item{ID: "d5ba79fc34d1fec8405a12a0525d63819d02906f", Type: rd.TypeR, Time: time.Date(2016, 12, 27, 10, 3, 4, 11, time.Local), Account: [2]string{"CH1234", "MO6677"}, Unit: 1.1, Amount: 100}, n: 0,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 100), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}, {ID: "MO6677", Type: "OK", Input: SubAcct(1, 100), Unit: SubUnit(1, 0), NUV: SubNuv(1, 0)}},
		},
		{ // del x & error
			ok: false, f: DelRecord,
			v: rd.Item{ID: "7", Type: rd.TypeX}, n: 0,
			a: []db.Account{},
		},
		{ // add i
			ok: true, f: AddRecord,
			v: rd.Item{ID: "8", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 100.5, Note: "hello", Member: "jj", Proj: "xx"}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 100.5), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
		},
		{ // chg i
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "cc35e99c2b4387b9e5db52782266475d84444947", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 100.5, Note: "hello", Member: "jj", Proj: "xx"}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 0), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
		},
		{ // chg i
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "cc35e99c2b4387b9e5db52782266475d84444947", Type: rd.TypeI, Time: time.Date(2016, 12, 25, 12, 3, 4, 11, time.Local), Account: [2]string{"MO6677", ""}, Unit: 10, Amount: 100.5, Note: "hello", Member: "jj", Proj: "xx"}, n: 1,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 100.5), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}, {ID: "MO6677", Type: "OK", Input: AddAcct(1, 100.5), Unit: AddUnit(1, 0), NUV: AddNuv(1, 0)}},
		},
		{ // add o
			ok: true, f: AddRecord,
			v: rd.Item{ID: "11", Type: rd.TypeO, Time: min, Account: [2]string{"CH1234", ""}, Unit: 5.5, Amount: 12.23, Member: "cc"}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 12.23), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}},
		},
		{ // chg o
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "e54141dcbe8218b324798f42be4f9039103e512b", Type: rd.TypeO, Time: time.Date(2017, 1, 5, 12, 3, 4, 11, time.Local), Account: [2]string{"CH1234", ""}, Unit: 5, Amount: 12, Member: "bb"}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 0.23), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
		},
		{ // add b
			ok: true, f: AddRecord,
			v: rd.Item{ID: "13", Type: rd.TypeB, Time: time.Date(2016, 12, 24, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 50, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}, n: 3,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 50), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
			d: []db.Debit{{ID: "630e62e41ed5e60983e5a75a3b16f813bd63b406", Name: "zhang3", Amount: 50, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		{ // add b
			ok: false, f: AddRecord,
			v: rd.Item{ID: "14", Type: rd.TypeB, Time: time.Date(2016, 12, 24, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"CH1234", ""}, Unit: 10, Amount: 50, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}, n: 3,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 0), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
			d: []db.Debit{{ID: "630e62e41ed5e60983e5a75a3b16f813bd63b406", Name: "zhang3", Amount: 50, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		{ // add l
			ok: true, f: AddRecord,
			v: rd.Item{ID: "15", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Member: "li4", Account: [2]string{"CH1234", ""}, Amount: 60}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 60), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}},
			d: []db.Debit{{ID: "9e1fcaf8eb8aab49860892e908ce4c5852b09976", Name: "li4", Amount: 60}},
		},
		{ // add l
			ok: false, f: AddRecord,
			v: rd.Item{ID: "16", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Account: [2]string{"CH1234", ""}, Amount: 60}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 0), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}},
			d: []db.Debit{{ID: "9e1fcaf8eb8aab49860892e908ce4c5852b09976", Name: "li4", Amount: 60}},
		},
		{ // chg b
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "ec6013a35aa2990b5f546e8862da8991066bd447", Type: rd.TypeB, Time: time.Date(2016, 12, 25, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"MO6677", ""}, Unit: 5, Amount: 10, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 50), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}, {ID: "MO6677", Type: "OK", Input: AddAcct(1, 10), Unit: AddUnit(1, 0), NUV: AddNuv(1, 0)}},
			d: []db.Debit{{ID: "630e62e41ed5e60983e5a75a3b16f813bd63b406", Name: "zhang3", Amount: 10, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}},
		},
		{ // chg b
			ok: false, f: ChgRecord,
			v: rd.Item{ID: "ec6013a35aa2990b5f546e8862da8991066bd447", Type: rd.TypeB, Time: time.Date(2016, 12, 25, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"MO6677", ""}, Unit: 5, Amount: 10, Deadline: time.Date(2017, 1, 1, 1, 0, 0, 0, time.Local)}, n: 4,
		},
		{ // del b
			ok: true, f: DelRecord,
			v: rd.Item{ID: "95b65a60dab4b1872fcf11db9da04ac3e914895c", Type: rd.TypeB, Time: time.Date(2016, 12, 25, 23, 18, 19, 0, time.Local), Member: "zhang3", Account: [2]string{"MO6677", ""}, Unit: 5, Amount: 10, Deadline: time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)}, n: 3,
			a: []db.Account{{ID: "MO6677", Type: "OK", Input: SubAcct(1, 10), Unit: SubUnit(1, 0), NUV: SubNuv(1, 0)}},
			d: []db.Debit{{ID: "630e62e41ed5e60983e5a75a3b16f813bd63b406", Note: "Del"}},
		},
		{ // chg l
			ok: true, f: ChgRecord,
			v: rd.Item{ID: "7bc80a2afbb7ce2962604380f1acd218eb02b38b", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Member: "li4", Account: [2]string{"CH1234", ""}, Amount: 70}, n: 3,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 10), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}},
			d: []db.Debit{{ID: "9e1fcaf8eb8aab49860892e908ce4c5852b09976", Name: "li4", Amount: 70}},
		},
		{ // del l
			ok: true, f: DelRecord,
			v: rd.Item{ID: "1f2e6aaf03a8f2c4804d40db7b5f3bc8b86c6547", Type: rd.TypeL, Time: time.Date(2016, 12, 26, 12, 3, 4, 0, time.Local), Member: "li4", Account: [2]string{"CH1234", ""}, Amount: 70}, n: 2,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 70), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
			d: []db.Debit{{ID: "9e1fcaf8eb8aab49860892e908ce4c5852b09976", Note: "Del"}},
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

	AddAcct := func(i int, v float64) float64 {
		testAcct[i] += v
		return testAcct[i]
	}

	SubAcct := func(i int, v float64) float64 {
		testAcct[i] -= v
		return testAcct[i]
	}

	AddUnit := func(i int, v float64) float64 {
		testUnit[i] += v
		return testUnit[i]
	}

	SubUnit := func(i int, v float64) float64 {
		testUnit[i] -= v
		if testUnit[i] < 0 {
			testUnit[i] = 0
		}
		return testUnit[i]
	}

	AddNuv := func(i int, v float64) float64 {
		testNuv[i] += v
		return testNuv[i]
	}

	SubNuv := func(i int, v float64) float64 {
		testNuv[i] -= v
		if testNuv[i] < 0 {
			testNuv[i] = 0
		}
		return testNuv[i]
	}

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
		// { // add a repeat
		// 	ok: false,
		// 	v:  Account{Account: db.Account{ID: "CD9007", Type: "A1", Unit: 0.55, NUV: 0.34, Class: "C1", Input: 100.55, Deadline: time.Date(2017, 2, 1, 2, 23, 12, 22, time.Local)}}, n: 1,
		// },
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
		{ // add b
			ok: true,
			v:  Account{Debit: db.Debit{Type: db.Borrow, Name: "zhang3", Amount: 16.55, Account: "CH1234", Note: "CC", Deadline: time.Date(2020, 2, 1, 2, 23, 12, 22, time.Local)}}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 16.55), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
			d: []db.Debit{{ID: "1bce89aa49d62a57d52ad4ec85d280a86af1dbdd", Name: "zhang3", Amount: 16.55, Account: "CH1234", Note: "CC", Deadline: time.Date(2020, 2, 1, 2, 23, 12, 22, time.Local)}},
		},
		{ // add b repeat
			ok: true,
			v:  Account{Debit: db.Debit{Type: db.Borrow, Name: "zhang3", Amount: 16.55, Account: "CH1234", Note: "CC", Deadline: time.Date(2020, 2, 1, 2, 23, 12, 22, time.Local)}}, n: 4,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 0), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}},
			d: []db.Debit{{ID: "1bce89aa49d62a57d52ad4ec85d280a86af1dbdd", Name: "zhang3", Amount: 16.55, Account: "CH1234", Note: "CC", Deadline: time.Date(2020, 2, 1, 2, 23, 12, 22, time.Local)}},
		},
		{ // chg b
			ok: true,
			v:  Account{Debit: db.Debit{Type: db.Borrow, Name: "zhang3", Amount: 10.55, Account: "MO6677", Note: "CC", Deadline: time.Date(2020, 2, 1, 2, 23, 12, 22, time.Local)}}, n: 5,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 16.55), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}, {ID: "MO6677", Type: "OK", Input: AddAcct(1, 10.55), Unit: AddUnit(1, 0), NUV: AddNuv(1, 0)}},
			d: []db.Debit{{ID: "1bce89aa49d62a57d52ad4ec85d280a86af1dbdd", Name: "zhang3", Amount: 10.55, Account: "MO6677", Note: "CC", Deadline: time.Date(2020, 2, 1, 2, 23, 12, 22, time.Local)}},
		},
		{ // del b
			ok: true,
			v:  Account{Debit: db.Debit{Type: db.Borrow, Name: "zhang3", Amount: 0, Account: "CH1234", Note: "CC", Deadline: time.Date(2020, 2, 1, 2, 23, 12, 22, time.Local)}}, n: 6,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 0), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}, {ID: "MO6677", Type: "OK", Input: SubAcct(1, 10.55), Unit: SubUnit(1, 0), NUV: SubNuv(1, 0)}},
			d: []db.Debit{{ID: "1bce89aa49d62a57d52ad4ec85d280a86af1dbdd", Note: "Del"}},
		},
		{ // add l
			ok: true,
			v:  Account{Debit: db.Debit{Type: db.Lend, Name: "li4", Amount: 50, Account: "CH1234"}}, n: 7,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 50), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}},
			d: []db.Debit{{ID: "9e1fcaf8eb8aab49860892e908ce4c5852b09976", Name: "li4", Amount: 50, Account: "CH1234"}},
		},
		{ // chg l
			ok: true,
			v:  Account{Debit: db.Debit{Type: db.Lend, Name: "li4", Amount: 90, Account: "CH1234"}}, n: 8,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: SubAcct(0, 40), Unit: SubUnit(0, 0), NUV: SubNuv(0, 0)}},
			d: []db.Debit{{ID: "9e1fcaf8eb8aab49860892e908ce4c5852b09976", Name: "li4", Amount: 90, Account: "CH1234"}},
		},
		{ // del l
			ok: true,
			v:  Account{Debit: db.Debit{Type: db.Lend, ID: "9e1fcaf8eb8aab49860892e908ce4c5852b09976"}}, n: 9,
			a: []db.Account{{ID: "CH1234", Type: "FREE", Input: AddAcct(0, 90), Unit: AddUnit(0, 0), NUV: AddNuv(0, 0)}},
			d: []db.Debit{{ID: "9e1fcaf8eb8aab49860892e908ce4c5852b09976", Note: "Del"}},
		},
	} {
		err := tc.v.Update()
		if err != nil { // record func
			if !tc.ok {
				continue
			}
			t.Fatalf("%d:%v", i, err)
		}
		if !tc.ok {
			t.Fatalf("%d, should be error", i)
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
				if !tc.ok || d.Note == "Del" {
					continue
				}
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
}

func TestAccounts(t *testing.T) { // TODO: not over
	{
		l1 := db.Debit{Name: "zhang3", Amount: 16}
		if err := db.GetLend().Add(l1); err != nil {
			t.Fatal(err)
		}
		b1 := db.Debit{Name: "li4", Amount: 55, Note: "Hello"}
		if err := db.GetBorrow().Add(b1); err != nil {
			t.Fatal(err)
		}
	}
	a := []Account{
		Account{Account: db.Account{ID: "CH1234", Type: "FREE", Input: testAcct[0], Unit: testUnit[0], NUV: testNuv[0]}},
		Account{Account: db.Account{ID: "MO6677", Type: "OK", Input: testAcct[1], Unit: testUnit[1], NUV: testNuv[1]}},
		Account{Debit: db.Debit{ID: "b94d706fe7f3369a446fe3ab68c73b3f89ca8e13", Type: db.Lend, Amount: 16, Name: "zhang3"}},
		Account{Debit: db.Debit{ID: "d7a5a4f649f2209c254777d8cab7e490166dda2f", Type: db.Borrow, Amount: 55, Name: "li4", Note: "Hello"}},
	}
	v, err := Accounts()
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", string(b))
	// clear time info
	var A []Account
	for _, vv := range v {
		vv.Account.Time = time.Time{}
		vv.Debit.Time = time.Time{}
		A = append(A, vv)
	}
	if !reflect.DeepEqual(a, A) {
		t.Fatalf("%+v\n%+v", a, A)
	}
}

func TestAccountJson(t *testing.T) {
	for i, a := range []Account{
		Account{Account: db.Account{Time: time.Date(2016, 2, 1, 2, 23, 12, 0, time.Local), ID: "CD9007", Type: "A1", Unit: 0.55, NUV: 0.34, Class: "C1", Input: 100.55, Deadline: time.Date(2017, 2, 1, 2, 23, 12, 0, time.Local)}},
		Account{},
		Account{Debit: db.Debit{Type: db.Borrow, Name: "zhang3", Amount: 16.55, Account: "CH1234", Note: "CC", Deadline: time.Date(2020, 2, 1, 2, 23, 12, 0, time.Local)}},
		Account{Debit: db.Debit{Type: db.Lend, Name: "zhang3", Amount: 16.55}},
	} {
		b, err := json.Marshal(a)
		if err != nil {
			t.Fatalf("%d:%v", i, err)
		}
		// fmt.Println(string(b))
		var v Account
		if err := json.Unmarshal(b, &v); err != nil {
			t.Fatal(err)
		}
		// fmt.Println(v)
		if !reflect.DeepEqual(a, v) {
			t.Fatalf("%d: %+v,%+v", i, v, a)
		}
	}
}

func TestRemove(t *testing.T) {
	os.RemoveAll("db")
}
