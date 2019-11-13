package controllers

import (
	"encoding/json"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wishlily/dashboard/server/account"
	db "github.com/wishlily/dashboard/server/account/database"
)

func TestPre(t *testing.T) {
	if err := Finance.Init(); err != nil {
		t.Fatal(err)
	}
}

func TestRecordsOk(t *testing.T) {
	/*
		{
			"key":   "xxx",
			"type":  "xxx",
			"time":  "yyyy-mm-dd hh:mm:ss",
			"amount": 0.0,
			"account":["xxx",...],
			"unit":   0.0,                    // omitempty
			"class":  ["xxx",...],            // omitempty
			"member": "xxx",                  // omitempty
			"proj":   "xxx",                  // omitempty
			"note":   "hello",                // omitempty
			"deadline":"yyyy-mm-dd hh:mm:ss", // omitempty
		}
	*/
	const URL = "/api/finance/record"

	db.GetAccount().Add(db.Account{ID: "AB1234", Type: "FREE", Input: 100})
	defer db.GetAccount().Del(db.Account{ID: "AB1234"})

	for i, tc := range []struct {
		f      func(string, string, *gin.Engine) (int, []byte)
		param  string
		result string
	}{
		{
			f: post,
			param: `{
				"type": "add",
				"data": {
					"type": "I",
					"time": "2018-03-20 12:23:54",
					"amount": 12.45,
					"account": ["AB1234"]
				}
			}`,
			result: `{"message":"ok"}`,
		},
		{
			f:      get,
			param:  "start=" + url.QueryEscape("2018-03-19 00:12:34") + "&end=" + url.QueryEscape("2018-03-20 13:10:00"),
			result: `[{"key":"c9c91795ea21b9966c010cea09b0e0fef1e10dd8","amount":12.45,"type":"I","time":"2018-03-20 12:23:54","account":["AB1234"]}]`,
		},
		{
			f: post,
			param: `{
				"type": "chg",
				"data": {
					"key": "c9c91795ea21b9966c010cea09b0e0fef1e10dd8",
					"type": "O",
					"time": "2018-03-19 12:23:54",
					"amount": 20,
					"account": ["AB1234"],
					"unit":   10,
					"class": ["1", "2"],
					"member": "li4",
					"proj":   "he",
					"note":   "hello",
					"deadline":"2018-05-01 01:20:30"
				}
			}`,
			result: `{"message":"ok"}`,
		},
		{
			f:      get,
			param:  "start=" + url.QueryEscape("2018-03-19 00:12:34") + "&end=" + url.QueryEscape("2018-03-20 13:10:00"),
			result: `[{"key":"89de8e7f0e207975ef5ce9c4483a5722a4ee9544","amount":20,"member":"li4","proj":"he","unit":10,"note":"hello","type":"O","time":"2018-03-19 12:23:54","class":["1","2"],"account":["AB1234"],"deadline":"2018-05-01 01:20:30"}]`,
		},
		{
			f: post,
			param: `{
				"type": "del",
				"data": {
					"key": "89de8e7f0e207975ef5ce9c4483a5722a4ee9544",
					"type": "O",
					"time": "2018-03-19 12:23:54",
					"amount": 20,
					"account": ["AB1234"]
				}
			}`,
			result: `{"message":"ok"}`,
		},
		{
			f:      get,
			param:  "start=" + url.QueryEscape("2018-03-19 00:12:34") + "&end=" + url.QueryEscape("2018-03-20 13:10:00"),
			result: `[]`,
		},
	} {
		code, body := tc.f(URL, tc.param, router)
		result := string(body)
		if result != tc.result {
			t.Fatalf("%d: %d, %v", i, code, result)
		}
	}

	acct, _ := db.GetAccount().Sel("AB1234")
	if acct.Input != 100 {
		t.Fatal()
	}
}

func TestRecordsErr(t *testing.T) {
	const URL = "/api/finance/record"
	for i, tc := range []struct {
		f     func(string, string, *gin.Engine) (int, []byte)
		param string
		code  int
	}{
		{
			f: post,
			param: `{
				"type": "add",
				"data": {
					"type": "I",
					"time": "2018-03-20 12:23:54",
					"amount": 12.45,
					"account": ["AB1234"]
				}
			}`,
			code: 500,
		},
		{
			f: post,
			param: `{
				"type": "chg",
				"data": {
					"type": "I",
					"time": "2018-03-20 12:23:54",
					"amount": 12.45,
					"account": ["AB1234"]
				}
			}`,
			code: 500,
		},
		{
			f: post,
			param: `{
				"type": "del",
				"data": {
					"type": "I",
					"time": "2018-03-20 12:23:54",
					"amount": 12.45,
					"account": ["AB1234"]
				}
			}`,
			code: 500,
		},
		{
			f: post,
			param: `{
				"type": "chg",
				"data": {
				}
			}`,
			code: 400,
		},
		{
			f: post,
			param: `{
				"type": "null",
				"data": {
					"type": "I",
					"time": "2018-03-20 12:23:54",
					"amount": 12.45,
					"account": ["AB1234"]
				}
			}`,
			code: 400,
		},
		{
			f:     get,
			param: "hh",
			code:  400,
		},
		{
			f:     get,
			param: "start=" + url.QueryEscape("2019-03-19 00:12:34") + "&end=" + url.QueryEscape("2018-03-20 13:10:00"),
			code:  500,
		},
	} {
		code, body := tc.f(URL, tc.param, router)
		result := string(body)
		if code != tc.code {
			t.Fatalf("%d: %d, %v", i, code, result)
		}
	}
}

func TestAccountsOk(t *testing.T) {
	parse := func(data []byte) ([]account.Account, error) {
		var v, temp []account.Account
		if err := json.Unmarshal(data, &temp); err != nil {
			return nil, err
		}
		for _, t := range temp {
			t.Account.Time = time.Time{} // not cmp
			t.Debit.Time = time.Time{}   // not cmp
			v = append(v, t)
		}
		return v, nil
	}
	/*
		{
		    "time":    "yyyy-mm-dd hh:mm:ss",
		    "id":      "xxx",
		    "type":    "xxx",
		    "amount":  0.0,
		    "unit":    0.0,                   // omitempty
		    "nuv":     0.0,                   // omitempty
		    "class":   "xxx",                 // omitempty
		    "deadline":"yyyy-mm-dd hh:mm:ss", // omitempty
		    "member":  "xxx",                 // omitempty
		    "account": "xxx",                 // omitempty
		    "note":    "xxx",                 // omitempty
		}
	*/
	const URL = "/api/finance/account"

	for i, tc := range []struct {
		f      func(string, string, *gin.Engine) (int, []byte)
		param  string
		result string
	}{
		{
			f: post,
			param: `{
				"type": "add",
				"data": {
					"id": "AB1234",
					"type": "AA",
					"amount": 100
				}
			}`,
			result: `{"message":"ok"}`,
		},
		{
			f:      get,
			param:  "",
			result: `[{"id":"AB1234","type":"AA","amount":100,"time":"2019-11-13 13:53:10"}]`,
		},
		{
			f: post,
			param: `{
				"type": "chg",
				"data": {
					"id":      "AB1234",
					"type":    "BB",
					"amount":  250.5,
					"unit":    18.5,
					"nuv":     2.5,
					"class":   "ABC",
					"deadline":"2019-12-21 12:34:56",
					"member":  "abc",
					"account": "def",
					"note":    "xxx"
				}
			}`,
			result: `{"message":"ok"}`,
		},
		{
			f:      get,
			param:  "",
			result: `[{"id":"AB1234","type":"BB","unit":18.5,"nuv":2.5,"class":"ABC","amount":250.5,"time":"2019-11-13 13:58:17","deadline":"2019-12-21 12:34:56"}]`,
		},
		{
			f: post,
			param: `{
				"type": "add",
				"data": {
					"id": "CD9876",
					"type": "BB",
					"amount": 0
				}
			}`,
			result: `{"message":"ok"}`,
		},
		{
			f:      get,
			param:  "",
			result: `[{"id":"AB1234","type":"BB","unit":18.5,"nuv":2.5,"class":"ABC","amount":250.5,"time":"2019-11-13 14:00:55","deadline":"2019-12-21 12:34:56"},{"id":"CD9876","type":"BB","amount":0,"time":"2019-11-13 14:00:55"}]`,
		},
		{
			f:      get,
			param:  "list",
			result: `["AB1234","CD9876"]`,
		},
		{
			f: post,
			param: `{
				"type": "del",
				"data": {
					"id": "CD9876"
				}
			}`,
			result: `{"message":"ok"}`,
		},
		{
			f: post,
			param: `{
				"type": "del",
				"data": {
					"id": "AB1234"
				}
			}`,
			result: `{"message":"ok"}`,
		},
		{
			f:      get,
			param:  "",
			result: `[]`,
		},
	} {
		code, body := tc.f(URL, tc.param, router)
		v, err := parse(body)
		if err != nil {
			if string(body) == tc.result { // {"message":"ok"}
				continue
			}
			t.Fatalf("%d: %v, %v", i, err, string(body))
		}
		a, err := parse([]byte(tc.result))
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		if code != 200 || !reflect.DeepEqual(a, v) {
			t.Fatalf("%d: %d\n%v\n%v", i, code, tc.result, string(body))
		}
	}
}

func TestAccountsErr(t *testing.T) {
	const URL = "/api/finance/account"

	for i, tc := range []struct {
		f     func(string, string, *gin.Engine) (int, []byte)
		param string
		code  int
	}{
		{
			f: post,
			param: `{
				"type": "hh",
				"data": {
					"id": "AB1234",
					"type": "AA",
					"amount": 100
				}
			}`,
			code: 400,
		},
		{
			f: post,
			param: `{
				"data": {
				}
			}`,
			code: 400,
		},
		{
			f: post,
			param: `{
				"type": "del",
				"data": {
					"id": "AB1234",
					"type": "AA",
					"amount": 100
				}
			}`,
			code: 500,
		},
		{
			f: post,
			param: `{
				"type": "chg",
				"data": {
					"id": "",
					"type": "",
					"amount": 0
				}
			}`,
			code: 500,
		},
	} {
		code, body := tc.f(URL, tc.param, router)
		if code != tc.code {
			t.Fatalf("%d: %v, %v", i, code, string(body))
		}
	}
}

func TestRemove(t *testing.T) {
	os.RemoveAll("db")
}
