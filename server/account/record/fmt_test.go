package record

import (
	"io/ioutil"
	"os"
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

func TestFmtTime(t *testing.T) {
	var f format
	for i, tc := range []struct {
		v string
	}{
		{"2019-08-25 12:09:31"},
		{"1870-01-01 08:00:00"},
	} {
		f.Time = tc.v
		v := f.time()
		if v.Format("2006-01-02 15:04:05") != tc.v {
			t.Fatalf("%d: should be %v: %v", i, tc.v, v)
		}
	}
}

func TestReaderAll(t *testing.T) {
	const f = "test/t1.csv"
	r := reader{f}
	defer testCSVClr(t, f)

	for i, tc := range []struct {
		b []byte
		v []format
	}{
		{
			[]byte(`type,time,class,account,amount,note
1,2,3,4,5.6,7
a,b,c,d,-5.6,e`),
			[]format{
				format{"1", "2", "3", "4", 5.6, "7"},
				format{"a", "b", "c", "d", -5.6, "e"},
			},
		},
		{ // no header read none
			[]byte(`1,2,3,4,5.6,7
a,b,c,d,-5.6,e`),
			[]format{format{}},
		},
	} {
		testCSVSet(t, f, tc.b)
		v, err := r.all()
		if err != nil {
			t.Fatalf("%d: read all %v", i, err)
		}
		if !reflect.DeepEqual(v, tc.v) {
			t.Fatalf("%d: read should be %v: %v", i, tc.v, v)
		}
	}
}

func TestWriterHeader(t *testing.T) {
	const f = "test/t2.csv"
	f1, _ := os.Create(f)
	f1.Close()
	defer os.Remove(f)

	w := writer{f}
	if err := w.header(); err != nil {
		t.Fatal(err)
	}
	v, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	h := []byte(`type,time,class,account,amount,note
`)
	if !reflect.DeepEqual(h, v) {
		t.Fatalf("header: %s", string(v))
	}
}

func TestReaderIs(t *testing.T) {
	const f = "test/t3.csv"
	defer os.Remove(f)

	r := reader{f}
	v, err := r.all()
	if err != nil {
		t.Fatal(err)
	}
	if a := []format{}; !reflect.DeepEqual(a, v) {
		t.Fatal(v)
	}
}

func TestWriterAll(t *testing.T) {
	const (
		f    = "test/t4.csv"
		text = `type,time,class,account,amount,note
1,2,3,4,5.6,7
a,b,c,d,-5.6,e`
	)
	testCSVSet(t, f, []byte(text))
	defer testCSVClr(t, f)
	w := writer{f}
	r := reader{f}

	a := []format{
		format{"A", "B", "C", "D", 1.1, "E"},
		format{"f", "g", "h", "i", -2.2, "j"},
	}
	if err := w.all(a); err != nil {
		t.Fatal(err)
	}
	v, err := r.all()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, a) {
		t.Fatal(v)
	}
}

func TestWriterAppend(t *testing.T) {
	const (
		f = "test/t5.csv"
	)
	defer os.Remove(f)
	w := writer{f}
	r := reader{f}

	if err := w.append(format{}); err == nil {
		t.Fatal("append not exist file should be error")
	}

	r.all()
	l := []format{}
	for i, tc := range []struct {
		a format
	}{
		{format{"f", "g", "h", "i", 2.2, "j"}},
		{format{"A", "B", "C", "D", 1.1, "E"}},
		{format{"1", "2", "3", "4", 5.6, "7"}},
	} {
		l = append(l, tc.a)
		if err := w.append(tc.a); err != nil {
			t.Fatalf("%d: append %v", i, err)
		}
		v, err := r.all()
		if err != nil {
			t.Fatalf("%d: read %v", i, err)
		}
		if !reflect.DeepEqual(l, v) {
			t.Fatalf("%d: cmp %v", i, v)
		}
	}
}
