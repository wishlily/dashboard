package record

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func testDBSet(t *testing.T) {
	const (
		p  = "test/"
		n1 = "2016.csv"
		v1 = `type,time,class,account,amount,note
R,2016-05-15 00:00:00,AA#BB,XX1#YY1,12.5,hello
I,2016-01-13 00:00:00,cc#ss,BB1,30,
`
		n2 = "2017.csv"
		v2 = `type,time,class,account,amount,note
O,2017-02-21 00:00:00,DD#zz,fff,89.34,
X,2017-03-20 00:00:00,gg#kl,adc,21,
`
	)
	testCSVSet(t, p+n1, []byte(v1))
	testCSVSet(t, p+n2, []byte(v2))
}

func testDBClr(t *testing.T) {
	const (
		p  = "test/"
		n1 = "2016.csv"
		n2 = "2017.csv"
	)
	testCSVClr(t, p+n1)
	testCSVClr(t, p+n2)
}

func TestDBPath(t *testing.T) {
	db := database{}
	for i, tc := range []struct {
		a string
		v string
	}{
		// {"./", "."},
		// {"", "."},
		// {"../", ".."},
		{"/root", "/root"},
		{"/root/", "/root"},
	} {
		if err := db.setPath(tc.a); err != nil {
			t.Fatalf("%d:%v", i, err)
		}
		if tc.v != db.path {
			t.Fatalf("%d: set path should be %v: %v", i, tc.v, db.path)
		}
	}
}

func TestDBCSV(t *testing.T) {
	db := database{path: "/root"}
	for i, tc := range []struct {
		y int
		v string
	}{
		{0, "/root/0000.csv"},
		{30, "/root/0030.csv"},
		{2016, "/root/2016.csv"},
	} {
		v := db.csv(tc.y)
		if v != tc.v {
			t.Fatalf("%d: should be %v: %v", i, tc.v, v)
		}
	}
}

func TestDBLoad(t *testing.T) {
	testDBSet(t)
	defer testDBClr(t)

	db := database{min: 2016, max: 2017, path: "test"}
	db.load()
	num := 0
	db.buf.Range(func(key interface{}, val interface{}) bool {
		num++
		fmt.Printf("%v -> %v\n", key, val)
		return true
	})
	if num != 4 {
		t.Fatalf("num should be 4 : %d", num)
	}
}

func TestDBPub(t *testing.T) {
	testDBSet(t)
	defer testDBClr(t)

	if err := SetPath("test"); err != nil {
		t.Fatal(err)
	}
	// Get
	its, err := Get("2016-05-15 00:00:00", "2017-02-21 00:00:00")
	if err != nil {
		t.Fatal(err)
	}
	if len(its) != 2 {
		t.Fatalf("Get should be 2 : %v", len(its))
	}
	fmt.Println(its)
	bak := its
	bak[0].ID = ""
	bak[1].ID = ""
	// Add
	t1, _ := time.ParseInLocation(timeFMT, "2016-06-01 12:23:54", time.Local)
	data := Item{Type: TypeL, Time: t1, Amount: 16, Note: "Add"}
	if err := Add(data); err != nil {
		t.Fatal(err)
	}
	its, err = Get("2016-05-15 00:00:00", "2017-02-21 00:00:00")
	if err != nil {
		t.Fatal(err)
	}
	if len(its) != 3 || its[1].Note != "Add" {
		t.Fatalf("Add should be 3 : %v", len(its))
	}
	fmt.Println(its)
	data.ID = its[1].ID // update ID
	// Chg
	data.Note = "Chg"
	t2, _ := time.ParseInLocation(timeFMT, "2017-01-01 12:23:54", time.Local)
	data.Time = t2
	if err := Chg(data); err != nil {
		t.Fatal(err)
	}
	// if !reflect.DeepEqual(old, its[1]) {
	// 	t.Fatalf("Chg old item : %v", old)
	// }
	its, err = Get("2016-05-15 00:00:00", "2017-02-21 00:00:00")
	if err != nil {
		t.Fatal(err)
	}
	if len(its) != 3 || its[1].Note != "Chg" || !its[1].Time.Equal(t2) {
		t.Fatalf("Chg should be 3 : %v", len(its))
	}
	fmt.Println(its)
	data.ID = its[1].ID // update ID
	// Sel
	if v, err := Sel(data); err != nil || !reflect.DeepEqual(v, its[1]) {
		t.Fatalf("%v:%v", err, v)
	}
	// Del
	if err := Del(data); err != nil {
		t.Fatal(err)
	}
	if _, err := Sel(data); err == nil {
		t.Fatal("Should be error")
	}
	its, err = Get("2016-05-15 00:00:00", "2017-02-21 00:00:00")
	if err != nil {
		t.Fatal(err)
	}
	if len(its) != 2 {
		t.Fatalf("Del should be 2 : %v", len(its))
	}
	fmt.Println(its)
	its[0].ID = ""
	its[1].ID = ""
	if !reflect.DeepEqual(bak, its) {
		t.Fatalf("Del : %v", its)
	}
}

func TestDBErr(t *testing.T) {
	data := Item{Time: time.Date(1999, 1, 1, 1, 1, 1, 1, time.Local)}
	if err := Add(data); err == nil {
		t.Fatalf("Add less 2000 year not support")
	}
	data.ID = "hello"
	if err := Del(data); err == nil {
		t.Fatalf("Del should be not found")
	}
	if err := Chg(data); err == nil {
		t.Fatalf("Chg should be not found")
	}
	db := database{}
	path := db.csv(2013)
	if path != "./2013.csv" {
		t.Fatalf("Get csv name %v", path)
	}
	if err := db.setRange(time.Unix(2019, 0), time.Unix(2016, 0)); err == nil {
		t.Fatalf("Need start < end")
	}
	if err := db.load(); err == nil {
		t.Fatalf("Need first set range")
	}
}
