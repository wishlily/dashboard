package database

import (
	"fmt"
	"time"
)

const (
	tableAccountName   = "account"
	tableAccountIDName = "id"
)

// Account database store account info
type Account struct {
	Time     int64   `db:"time" type:"INTEGER(64)"`
	ID       string  `db:"id" type:"VARCHAR(20)"`
	Type     string  `db:"type" type:"VARCHAR(3)"`
	Unit     float64 `db:"unit" type:"DECIMAL(32,3)"`
	NUV      float64 `db:"nuv" type:"DECIMAL(8,3)"` // net unit value
	Class    string  `db:"class" type:"VARCHAR(5)"`
	Input    float64 `db:"input" type:"DECIMAL(32,3)"`
	Deadline int64   `db:"deadline" type:"INTEGER(64)"`
}

type accountTable struct {
	table
}

func newAccountTable(url string) (accountTable, error) {
	db, err := newTable(url, tableAccountName, Account{})
	if err != nil {
		return accountTable{}, err
	}
	return accountTable{db}, nil
}

func (a accountTable) get() ([]Account, error) {
	rows, err := a.table.get()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accs := []Account{}
	for rows.Next() {
		var v Account
		err = rows.StructScan(&v)
		if err != nil {
			return nil, err
		}
		accs = append(accs, v)
	}
	return accs, nil
}

func (a accountTable) sel(id string) (Account, error) {
	v := Account{ID: id}
	rows, err := a.table.sel(tableAccountIDName, v)
	if err != nil {
		return v, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&v)
		if err != nil {
			return v, err
		}
		return v, nil
	}
	return v, fmt.Errorf("Not found id:%v in account table", id)
}

func (a accountTable) add(data Account) error {
	data.Time = time.Now().Unix()
	if _, err := a.sel(data.ID); err == nil {
		return fmt.Errorf("Add should be unique id : %v", data.ID)
	}
	return a.table.insert(data)
}

func (a accountTable) chg(data Account) error {
	data.Time = time.Now().Unix()
	return a.table.update(tableAccountIDName, data)
}

func (a accountTable) del(data Account) error {
	return a.table.delete(tableAccountIDName, data)
}
