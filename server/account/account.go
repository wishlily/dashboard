package account

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
	db "github.com/wishlily/dashboard/server/account/database"
	rd "github.com/wishlily/dashboard/server/account/record"
)

const (
	timeFMT = "2006-01-02 15:04:05"
)

// Init account config by viper
func Init() error {
	path := viper.GetString("csv.path")
	if len(path) == 0 {
		path = "db"
	}
	if err := rd.SetPath(path); err != nil {
		return err
	}
	url := viper.GetString("db.url")
	if len(url) == 0 {
		url = "db/database.db"
	}
	if err := db.SetURL(url); err != nil {
		return err
	}
	return nil
}

// Records get start (xxxx-xx-xx xx:xx:xx) to end (xxxx-xx-xx xx:xx:xx) records
func Records(start, end string) ([]Record, error) {
	its, err := rd.Get(start, end)
	if err != nil {
		return nil, err
	}
	records := []Record{}
	for _, v := range its {
		records = append(records, Record{v})
	}
	return records, nil
}

// Record one csv item
type Record struct {
	rd.Item
}

// Add one record
func (r Record) Add() error {
	switch r.Type {
	case rd.TypeI:
		return r.io(rd.Add, r.add)
	case rd.TypeO:
		return r.io(rd.Add, r.sub)
	case rd.TypeR:
		return r.r(rd.Add, r.sub, r.add)
	case rd.TypeB:
		return r.debit(rd.Add, db.GetBorrow().Add, r.add)
	case rd.TypeL:
		return r.debit(rd.Add, db.GetLend().Add, r.sub)
	}
	return fmt.Errorf("Can't support %v type record add", r.Type)
}

// Del one record, debit item Member&Note&Deadline can't change
func (r Record) Del() error {
	switch r.Type {
	case rd.TypeI:
		return r.io(rd.Del, r.sub)
	case rd.TypeO:
		return r.io(rd.Del, r.add)
	case rd.TypeR:
		return r.r(rd.Del, r.add, r.sub)
	case rd.TypeB:
		return r.debit(rd.Del, db.GetBorrow().Del, r.sub)
	case rd.TypeL:
		return r.debit(rd.Del, db.GetLend().Del, r.add)
	}
	return fmt.Errorf("Can't support %v type record add", r.Type)
}

// Chg one record, debit item Member&Note&Deadline can't change
func (r Record) Chg() error {
	v, err := rd.Sel(r.Item)
	if err != nil {
		return err
	}
	old := Record{v}
	if err := old.Del(); err != nil {
		return err
	}
	return r.Add()
}

func (r Record) add(acct *db.Account) {
	acct.Unit += r.Unit
	acct.Input += r.Amount
}

func (r Record) sub(acct *db.Account) {
	acct.Unit -= r.Unit
	if acct.Unit < 0 {
		acct.Unit = 0 // Unit not zero
	}
	acct.Input -= r.Amount
}

// f : rd.Add rd.Del rd.Chg
func (r Record) io(f func(rd.Item) error, opt func(acct *db.Account)) error {
	id := r.Account[rd.AccountM]
	acct, err := db.GetAccount().Sel(id)
	if err != nil {
		return err
	}
	opt(&acct)
	if err := f(r.Item); err != nil {
		return err
	}
	return db.GetAccount().Chg(acct)
}

func (r Record) r(f func(rd.Item) error, optM func(acct *db.Account), optS func(acct *db.Account)) error {
	idM := r.Account[rd.AccountM]
	acctM, err := db.GetAccount().Sel(idM)
	if err != nil {
		return err
	}
	idS := r.Account[rd.AccountS]
	acctS, err := db.GetAccount().Sel(idS)
	optM(&acctM)
	optS(&acctS)
	if err := f(r.Item); err != nil {
		return err
	}
	if err := db.GetAccount().Chg(acctM); err != nil {
		return err
	}
	return db.GetAccount().Chg(acctS)
}

func (r Record) debit(f func(rd.Item) error, fd func(db.Debit) error, opt func(acct *db.Account)) error {
	member := r.Member
	if len(member) == 0 {
		return fmt.Errorf("Should be have member")
	}
	id := r.Account[rd.AccountM]
	acct, err := db.GetAccount().Sel(id)
	if err != nil {
		return err
	}
	opt(&acct)
	data := db.Debit{Name: member, Amount: r.Amount, Note: r.Note, Deadline: r.Deadline}
	if err := fd(data); err != nil { // id by Name&Note&Deadline
		return err
	}
	if err := f(r.Item); err != nil {
		return err
	}
	return db.GetAccount().Chg(acct)
}

// MarshalJSON json encode
func (r Record) MarshalJSON() ([]byte, error) {
	type Item rd.Item
	item := &struct {
		Item
		Type     string   `json:"type"`
		Time     string   `json:"time"`
		Class    []string `json:"class,omitempty"`
		Account  []string `json:"account"`
		Deadline string   `json:"deadline,omitempty"`
	}{
		Item: (Item)(r.Item),
		Type: r.Type.String(),
		Time: r.Time.Format(timeFMT),
	}
	if r.Deadline.Unix() >= r.Time.Unix() { // is valid ?
		item.Deadline = r.Deadline.Format(timeFMT)
	}
	for i, v := range r.Class {
		if len(v) > 0 {
			item.Class = r.Class[:i+1]
		}
	}
	for i, v := range r.Account {
		if len(v) > 0 {
			item.Account = r.Account[:i+1]
		}
	}
	return json.Marshal(item)
}

// UnmarshalJSON json decode
func (r *Record) UnmarshalJSON(b []byte) error {
	type Item rd.Item
	item := &struct {
		*Item
		Type     string `json:"type"`
		Time     string `json:"time"`
		Deadline string `json:"deadline,omitempty"`
	}{
		Item: (*Item)(&r.Item),
	}
	err := json.Unmarshal(b, &item)
	if err != nil {
		return err
	}
	r.Type = rd.ParseType(item.Type)
	r.Time, err = time.ParseInLocation(timeFMT, item.Time, time.Local)
	if err != nil {
		return err
	}
	if len(item.Deadline) > 0 { // is have ?
		r.Deadline, err = time.ParseInLocation(timeFMT, item.Deadline, time.Local)
		if err != nil {
			return err
		}
	}
	return nil
}

// Account database data
type Account struct {
	Account db.Account
	Debit   db.Debit
}

// Update account/debit data
func (a Account) Update() error {
	if a.isAccount() {
		return a.updateAcct()
	}
	if a.isDebit() {
		// 	if a.Debit.Type == db.Borrow {
		// 		return
		// 	}
		// 	if a.Debit.Type == db.Lend {
		// 		return
		// 	}
	}
	return fmt.Errorf("No data in Account")
}

func (a Account) isAccount() bool {
	var v db.Account
	if a.Account != v { // empty
		return true
	}
	return false
}

func (a Account) isDebit() bool {
	var v db.Debit
	if a.Debit != v { // empty
		return true
	}
	return false
}

func (a Account) account() (db.Account, error) {
	old, err := db.GetAccount().Sel(a.Account.ID)
	if err != nil { // not found
		return old, db.GetAccount().Add(a.Account)
	}
	var empty db.Account
	acct := a.Account
	acct.ID = ""
	if acct == empty {
		return old, db.GetAccount().Del(a.Account)
	}
	return old, db.GetAccount().Chg(a.Account)
}

func (a Account) updateAcct() error {
	old, err := a.account()
	if err != nil {
		return err
	}
	// add record
	r := rd.Item{
		Type:     rd.TypeX,
		Time:     time.Now(),
		Amount:   a.Account.Input - old.Input,
		Unit:     a.Account.Unit - old.Unit,
		Deadline: a.Account.Deadline,
	}
	r.Account[rd.AccountM] = a.Account.ID
	r.Class[rd.ClassM] = a.Account.Type
	r.Class[rd.ClassS] = a.Account.Class
	if a.Account.NUV != 0 {
		r.Note = fmt.Sprintf("%v", a.Account.NUV)
	}
	if r.Amount == 0 && r.Unit == 0 {
		return nil // no add one record
	}
	return rd.Add(r)
}

func (a Account) debit() error {
	// TODO:
	debit, err:=
}

func (a Account) updebit(f func(db.Debit) error, opt func(*db.Account, float64)) error {
	// add/chg/del db
	if err := f(a.Debit); err != nil {
		return err
	}
	// account
	acct, err := db.GetAccount().Sel(a.Debit.Account)
	if err != nil {
		return err
	}
	if a.Debit.Type == db.Borrow {
		acct.Input += a.Debit.Amount
	} else {
		acct.Input -= a.Debit.Amount
	}
	// add record
	r := rd.Item{
		Type:     rd.TypeX,
		Time:     time.Now(),
		Amount:   a.Debit.Amount,
		Deadline: a.Debit.Deadline,
		Note:     a.Debit.Note,
		Member:   a.Debit.Name,
	}
	r.Account[rd.AccountM] = a.Debit.Account
	if err := rd.Add(r); err != nil {
		return err
	}
	return db.GetAccount().Chg(acct)
}
