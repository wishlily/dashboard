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
		"":  "xxx",
	}
	var it Item
	v := it.note(notes)
	if v != "xxx#a#1#b#2" &&
		v != "xxx#b#2#a#1" {
		t.Fatal(v)
	}

	notes = map[string]string{
		"a": "1",
		"b": "2",
	}
	v = it.note(notes)
	if v != "#a#1#b#2" &&
		v != "#b#2#a#1" {
		t.Fatal(v)
	}
}

func TestItemParseItem(t *testing.T) {
	v := format{
		Type: "O", Time: "2016-04-21 12:03:21",
		Class: "x1#y1", Account: "AA#BB", Amount: 12.12,
		Note: "xxx#member#mm#proj#pp#unit#23#deadline#2019-08-08 12:00:00#hello#hello",
	}
	item := parseItem(v, 2014, 0)
	if len(item.ID) == 0 || item.ID == "0" {
		t.Fatalf("ID not be 0")
	}
	item.ID = ""
	want := Item{csv: v, Type: TypeO, Time: time.Unix(1461211401, 0),
		Class: [classN]string{"x1", "y1"}, Amount: 12.12, Account: [accountN]string{"AA", "BB"},
		Member: "mm", Proj: "pp", Unit: 23, Deadline: time.Unix(1565236800, 0),
		Note: "xxx#hello#hello"}
	if !reflect.DeepEqual(item, want) {
		t.Fatalf("should be %v:%v", want, item)
	}
}

func TestItemUpdate(t *testing.T) {
	data := format{
		Type: "O", Time: "2016-04-21 12:03:21",
		Class: "x1#y1", Account: "AA#BB", Amount: 12.12,
		Note: "xxx#member#mm#proj#pp#unit#23#deadline#2019-08-08 12:00:00#hello#hello",
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
