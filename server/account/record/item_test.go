package record

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestItemNote(t *testing.T) {
	notes := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"":  "xxx",
	}
	var it Item
	v := it.note(notes)
	if v != "xxx#a#1#b#2#c#3" {
		t.Fatal(v)
	}

	notes = map[string]string{
		"a":  "1",
		"b":  "2",
		"_C": "3",
	}
	v = it.note(notes)
	if v != "#_C#3#a#1#b#2" {
		t.Fatal(v)
	}
}

func TestItemParseItem(t *testing.T) {
	v := format{
		Type: "O", Time: "2016-04-21 12:03:21",
		Class: "x1#y1", Account: "AA#BB", Amount: 12.12,
		Note: "xxx#member#mm#proj#pp#unit#23#deadline#2019-08-08 12:00:00#nuv#2.4#hello#hello",
	}
	item := parseItem(v, 2014, 0)
	if len(item.ID) == 0 || item.ID == "0" {
		t.Fatalf("ID not be 0")
	}
	item.ID = ""
	want := Item{csv: v, Type: TypeO, Time: time.Date(2016, 4, 21, 12, 3, 21, 0, time.Local),
		Class: [classN]string{"x1", "y1"}, Amount: 12.12, Account: [accountN]string{"AA", "BB"},
		Member: "mm", Proj: "pp", Unit: 23, NUV: 2.4, Deadline: time.Date(2019, 8, 8, 12, 0, 0, 0, time.Local),
		Note: "xxx#hello#hello"}
	if !reflect.DeepEqual(item, want) {
		t.Fatalf("should be %v:%v", want, item)
	}
}

func TestItemUpdate(t *testing.T) {
	data := format{
		Type: "O", Time: "2016-04-21 12:03:21",
		Class: "x1#y1", Account: "AA#BB", Amount: 12.12,
		Note: "xxx#member#mm#proj#pp#unit#23#deadline#2019-08-08 12:00:00#hello#hello#nuv#2.4",
	}
	it := parseItem(data, 2015, 0)
	it.csv = format{}
	if reflect.DeepEqual(it.csv, data) {
		t.Fatalf("The csv need clean")
	}
	it.update()
	it.ID = "" // clean
	fmt.Println(it.csv)
	v := parseItem(it.csv, 2015, 0)
	v.ID = "" // clean
	if !reflect.DeepEqual(v, it) {
		t.Fatalf("update lost some info : %v", v)
	}
}

func TestItemsSort(t *testing.T) {
	data := items([]Item{
		Item{Time: time.Unix(124719212, 0)},
		Item{Time: time.Unix(124719212-33, 0)},
		Item{Time: time.Unix(124719212+33, 0)},
	})
	want := items([]Item{
		Item{Time: time.Unix(124719212-33, 0)},
		Item{Time: time.Unix(124719212, 0)},
		Item{Time: time.Unix(124719212+33, 0)},
	})
	sort.Sort(data)
	if !reflect.DeepEqual(data, want) {
		t.Fatalf("sort:%v", data)
	}
}
