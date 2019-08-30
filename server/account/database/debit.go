package database

import (
	"fmt"
	"time"
)

const (
	tableBorrowName  = "borrow"
	tableLendName    = "lend"
	tableDebitIDName = "id"
)

// Debit borrow|lend database
type Debit struct {
	Time     int64   `db:"time" type:"INTEGER(64)"`
	ID       string  `db:"id" type:"VARCHAR(32)"`
	Name     string  `db:"name" type:"VARCHAR(20)"`
	Amount   float64 `db:"amount" type:"DECIMAL(32,3)"`
	Note     string  `db:"note" type:"VARCHAR(32)"`
	Deadline int64   `db:"deadline" type:"INTEGER(64)"`
}

type debitTable struct {
	table
}

func newBorrowTable(url string) (debitTable, error) {
	return newDebitTable(url, tableBorrowName)
}

func newLendTable(url string) (debitTable, error) {
	return newDebitTable(url, tableLendName)
}

func newDebitTable(url string, name string) (debitTable, error) {
	db, err := newTable(url, name, Debit{})
	if err != nil {
		return debitTable{}, err
	}
	return debitTable{db}, nil
}

func (d debitTable) get() ([]Debit, error) {
	rows, err := d.table.get()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	debs := []Debit{}
	for rows.Next() {
		var v Debit
		err = rows.StructScan(&v)
		if err != nil {
			return nil, err
		}
		debs = append(debs, v)
	}
	return debs, nil
}

func (d debitTable) sel(id string) (Debit, error) {
	v := Debit{ID: id}
	rows, err := d.table.sel(tableDebitIDName, v)
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
	return v, fmt.Errorf("Not found id:%v in %v table", id, d.table.name)
}

// gener one ID
func (d debitTable) id() string {
	return fmt.Sprintf("%v", time.Now().UnixNano())
}

func (d debitTable) add(data Debit) error {
	data.Time = time.Now().Unix()
	data.ID = d.id()
	if _, err := d.sel(data.ID); err == nil {
		return fmt.Errorf("Add should be unique id : %v", data.ID)
	}
	return d.table.insert(data)
}

func (d debitTable) chg(data Debit) error {
	data.Time = time.Now().Unix()
	return d.table.update(tableDebitIDName, data)
}

func (d debitTable) del(data Debit) error {
	return d.table.delete(tableDebitIDName, data)
}
