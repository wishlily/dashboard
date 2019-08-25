package record

import (
	"testing"
)

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

// func TestDBLoad(t *testing.T) {
// 	db := database{min: 2016, max: 2016}
// 	db.load()
// }

// const _TEST_CSV_FILENAME = "csv-test.tmp"

// func testPrepare(data []byte) error {
// 	if err := ioutil.WriteFile(_TEST_CSV_FILENAME, data, 0644); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func TestSplit(t *testing.T) {
// 	var m Bill
// 	m.Note = `其他内容#测试#内容 这样
// 	#实验# ad dsd########`
// 	note := m.split()
// 	out := map[string]string{
// 		"Default": "其他内容",
// 		"测试":      "内容 这样",
// 		"实验":      "ad dsd",
// 	}
// 	if !reflect.DeepEqual(out, note) {
// 		t.Fatalf("Split DeepEqual: \n%v\v%v", note, out)
// 	}
// }

// func TestSort(t *testing.T) {
// 	in := _Bills{
// 		Bill{Time: "2018-06-18 03:23:12"},
// 		Bill{Time: "2016-07-22 15:30:40"},
// 		Bill{Time: "2016-08-05 13:00:08"},
// 	}
// 	out := _Bills{
// 		Bill{Time: "2016-07-22 15:30:40"},
// 		Bill{Time: "2016-08-05 13:00:08"},
// 		Bill{Time: "2018-06-18 03:23:12"},
// 	}
// 	sort.Sort(in)
// 	if !reflect.DeepEqual(in, out) {
// 		t.Fatalf("Sort is not EXPECT: %v", in)
// 	}
// 	if !sort.IsSorted(in) {
// 		t.Fatalf("Sort IsSorted not TRUE")
// 	}
// }

// func TestReadAll(t *testing.T) {
// 	var in = []byte(`
// 交易类型,日期,类别,测试,子类,项目,成员,账户A,账户B,金额,份额,备注
// 收入,2018-05-23,活期,test,银行卡,,本人,4312,,20,1,
// 支出,2017-03-24,,test2,微信,,本人,2345,,12,1,
// `)
// 	out := _Bills{
// 		Bill{
// 			Time: "2017-03-24", Type: "支出", ClasSub: "微信", Member: "本人",
// 			Account: "2345", Amount: 12, Units: 1,
// 		},
// 		Bill{
// 			Time: "2018-05-23", Type: "收入", Class: "活期", ClasSub: "银行卡", Member: "本人",
// 			Account: "4312", Amount: 20, Units: 1,
// 		},
// 	}
// 	if err := testPrepare(in); err != nil {
// 		t.Fatalf("Prepare : %v", err)
// 	}

// 	c := NewCSV(_TEST_CSV_FILENAME)
// 	err := c.readAll()
// 	if err != nil {
// 		t.Fatalf("ReadAll %s: %v", _TEST_CSV_FILENAME, err)
// 	}
// 	for i, _ := range c.data {
// 		c.data[i].ID = ""
// 	}
// 	if !reflect.DeepEqual(c.data, out) {
// 		t.Fatalf("DeepEqual: \n%+v\n%+v", c.data, out)
// 	}
// 	os.Remove(_TEST_CSV_FILENAME)
// }

// func TestWrite(t *testing.T) {
// 	in := _Bills{
// 		Bill{
// 			Time: "2007-03-04", Type: "支出", ClasSub: "微信",
// 			Account: "2345",
// 		},
// 		Bill{
// 			Time: "2008-05-03", Type: "收入", Member: "本人",
// 			Account: "4312",
// 		},
// 	}

// 	c := NewCSV(_TEST_CSV_FILENAME)
// 	c.data = in
// 	if err := c.writeAll(); err != nil {
// 		t.Fatalf("WriteAll %s: %v", _TEST_CSV_FILENAME, err)
// 	}
// 	if err := c.readAll(); err != nil {
// 		t.Fatalf("ReadAll %s: %v", _TEST_CSV_FILENAME, err)
// 	}
// 	for i, _ := range c.data {
// 		c.data[i].ID = ""
// 	}
// 	if !reflect.DeepEqual(c.data, in) {
// 		t.Fatalf("DeepEqual: %v", c.data)
// 	}
// 	/*** Write Test ***/
// 	if err := c.write(in[0]); err != nil {
// 		t.Fatalf("Write: %v", err)
// 	}
// 	in = _Bills{
// 		Bill{
// 			Time: "2007-03-04", Type: "支出", ClasSub: "微信",
// 			Account: "2345",
// 		},
// 		Bill{
// 			Time: "2007-03-04", Type: "支出", ClasSub: "微信",
// 			Account: "2345",
// 		},
// 		Bill{
// 			Time: "2008-05-03", Type: "收入", Member: "本人",
// 			Account: "4312",
// 		},
// 	}
// 	if err := c.readAll(); err != nil {
// 		t.Fatalf("Write ReadAll: %v", err)
// 	}
// 	for i, _ := range c.data {
// 		c.data[i].ID = ""
// 	}
// 	if !reflect.DeepEqual(c.data, in) {
// 		t.Fatalf("Write DeepEqual: %v", c.data)
// 	}
// 	os.Remove(_TEST_CSV_FILENAME)
// }

// func TestDel(t *testing.T) {
// 	c := NewCSV(_TEST_CSV_FILENAME)
// 	c.cache = make(map[string]Bill)
// 	if err := c.Del("1"); err != errcode.NFIND {
// 		t.Fatalf("Del should have ERROR: %v", errcode.NFIND)
// 	}
// 	c.cache["1"] = Bill{
// 		Time: "2007-03-04", Type: "支出", ClasSub: "微信",
// 		Account: "2345",
// 	}
// 	c.cache["2"] = Bill{
// 		Time: "2008-05-03", Type: "收入", Member: "本人",
// 		Account: "4312",
// 	}
// 	if err := c.Del("1"); err != nil {
// 		t.Fatalf("Del: %v", err)
// 	}
// 	in := c.data
// 	for i, _ := range in {
// 		in[i].ID = ""
// 	}
// 	out, err := c.Data()
// 	if err != nil {
// 		t.Fatalf("Data: %v", err)
// 	}
// 	for i, _ := range out {
// 		out[i].ID = ""
// 	}
// 	if !reflect.DeepEqual(out, []Bill(in)) {
// 		t.Fatalf("Del DeepEqual: \n%v\n%v", in, out)
// 	}
// 	os.Remove(_TEST_CSV_FILENAME)
// }
