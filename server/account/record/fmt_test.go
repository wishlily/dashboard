package record

import (
	"reflect"
	"testing"
)

func TestFmtSplit(t *testing.T) {
	var f format
	a := "hello#world 123#\ni am# robot #"
	s := []string{"hello", "world 123", "i am", "robot", ""}
	v := f.split(a)
	if !reflect.DeepEqual(v, s) {
		t.Fatalf("split should be %v : %v", s, v)
	}
}

func TestFmtNote(t *testing.T) {
	var f format
	f.Note = `其他内容#测试#内容 这样
	#实验# ad dsd########`
	v := f.note()
	s := map[string]string{
		"":   "其他内容",
		"测试": "内容 这样",
		"实验": "ad dsd",
	}
	if !reflect.DeepEqual(v, s) {
		t.Fatalf("note should be %v : %v", s, v)
	}
}

func TestFmtClass(t *testing.T) {
	var f format
	for i, tc := range []struct {
		s string
		v []string
	}{
		{"hello#world", []string{"hello", "world"}},
		{"hi", []string{"hi"}},
	} {
		f.Class = tc.s
		v := f.class()
		if !reflect.DeepEqual(v, tc.v) {
			t.Fatalf("%d: class should be %v : %v", i, tc.v, v)
		}
	}
}

func TestFmtAccount(t *testing.T) {
	var f format
	for i, tc := range []struct {
		s string
		v []string
	}{
		{"hello#world", []string{"hello", "world"}},
		{"hi", []string{"hi"}},
	} {
		f.Account = tc.s
		v := f.account()
		if !reflect.DeepEqual(v, tc.v) {
			t.Fatalf("%d: class should be %v : %v", i, tc.v, v)
		}
	}
}
