package database

import (
	"fmt"
	"time"
)

const (
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

// DebitTable database debit table
type DebitTable struct {
	table
}

func newDebitTable(url string, name string) (DebitTable, error) {
	db, err := newTable(url, name, Debit{})
	if err != nil {
		return DebitTable{}, err
	}
	return DebitTable{db}, nil
}

// Get all Debit
func (d DebitTable) Get() ([]Debit, error) {
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

// Sel one debit by ID
func (d DebitTable) Sel(id string) (Debit, error) {
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
func (d DebitTable) id() string {
	return fmt.Sprintf("%v", time.Now().UnixNano())
}

// Add one Debit
func (d DebitTable) Add(data Debit) error {
	data.Time = time.Now().Unix()
	data.ID = d.id()
	if _, err := d.Sel(data.ID); err == nil {
		return fmt.Errorf("Add should be unique id : %v", data.ID)
	}
	return d.table.insert(data)
}

// Chg one Debit data
func (d DebitTable) Chg(data Debit) error {
	data.Time = time.Now().Unix()
	return d.table.update(tableDebitIDName, data)
}

// Del one Debit by ID
func (d DebitTable) Del(data Debit) error {
	return d.table.delete(tableDebitIDName, data)
}
