package record

import (
	"reflect"
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
