package database

import (
	"crypto/sha1"
	"fmt"
	"time"
)

// DebitType Account func input type
type DebitType int

const (
	// Borrow Account Debit type
	Borrow DebitType = iota
	// Lend Account Debit type
	Lend
)

func (d DebitType) String() string {
	if d == Borrow {
		return "B"
	} else if d == Lend {
		return "L"
	}
	return "N"
}

const (
	tableDebitIDName = "id"
)

// Debit borrow|lend database
type Debit struct {
	Time     time.Time `db:"time" type:"TIMESTAMP"`
	Type     DebitType `db:"-"`
	ID       string    `db:"id" type:"VARCHAR(32)"`
	Name     string    `db:"name" type:"VARCHAR(20)"`
	Amount   float64   `db:"amount" type:"DECIMAL(32,3)"`
	Account  string    `db:"account" type:"VARCHAR(20)"`
	Note     string    `db:"note" type:"VARCHAR(32)"`
	Deadline time.Time `db:"deadline" type:"TIMESTAMP"`
}

func (d Debit) hash() string {
	ss := fmt.Sprintf("%v:%v:%v", d.Name, d.Note, d.Deadline.Format("2006 15:04:05"))
	sum := sha1.Sum([]byte(ss))
	return fmt.Sprintf("%x", sum)
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
func (d DebitTable) Sel(data Debit) (Debit, error) {
	if len(data.ID) == 0 {
		data.ID = data.hash()
	}
	rows, err := d.table.sel(tableDebitIDName, data)
	if err != nil {
		return Debit{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&data)
		if err != nil {
			return Debit{}, err
		}
		return data, nil
	}
	return Debit{}, fmt.Errorf("Not found id:%v in %v table", data.ID, d.table.name)
}

// Add one Debit
func (d DebitTable) Add(data Debit) error {
	data.Time = time.Now()
	data.ID = data.hash()
	if _, err := d.Sel(data); err == nil {
		return fmt.Errorf("Add should be unique id : %v", data.ID)
	}
	return d.table.insert(data)
}

// Chg one Debit data
func (d DebitTable) Chg(data Debit) error {
	if len(data.ID) == 0 {
		data.ID = data.hash()
	}
	data.Time = time.Now()
	return d.table.update(tableDebitIDName, data)
}

// Del one Debit by ID
func (d DebitTable) Del(data Debit) error {
	if len(data.ID) == 0 {
		data.ID = data.hash()
	}
	return d.table.delete(tableDebitIDName, data)
}
