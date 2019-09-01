package database

import (
	"fmt"
	"time"
)

const (
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

// AccountTable database table
type AccountTable struct {
	table
}

func newAccountTable(url string, name string) (AccountTable, error) {
	db, err := newTable(url, name, Account{})
	if err != nil {
		return AccountTable{}, err
	}
	return AccountTable{db}, nil
}

// Get all Account
func (a AccountTable) Get() ([]Account, error) {
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

// Sel get one Account by ID
func (a AccountTable) Sel(id string) (Account, error) {
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

// Add one Account
func (a AccountTable) Add(data Account) error {
	data.Time = time.Now().Unix()
	if _, err := a.Sel(data.ID); err == nil {
		return fmt.Errorf("Add should be unique id : %v", data.ID)
	}
	return a.table.insert(data)
}

// Chg change one Account data
func (a AccountTable) Chg(data Account) error {
	data.Time = time.Now().Unix()
	return a.table.update(tableAccountIDName, data)
}

// Del the Account by ID
func (a AccountTable) Del(data Account) error {
	return a.table.delete(tableAccountIDName, data)
}
