package record

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// ClassM : Item.Class[] master no -> Item.Class[ClassM]
	ClassM = iota
	// ClassS : Item.Class[] sub no -> Item.Class[ClassS]
	ClassS
	classN
)

const (
	// AccountM : Item.Account[] master no -> Item.Account[AccountM]
	AccountM = iota
	// AccountS : Item.Account[] sub no -> Item.Account[AccountS]
	AccountS
	accountN
)

// Item csv record
type Item struct {
	csv      format
	ID       string // CAN'T CHG
	Type     Type
	Time     time.Time
	Class    [classN]string
	Amount   float64
	Account  [accountN]string
	Member   string
	Proj     string // project & id
	Unit     float64
	Deadline time.Time
	Note     string
}

func (it Item) id(year, num int) string {
	return fmt.Sprintf("%04d%v%d", year, time.Now().Unix(), num)
}

func (it Item) note(notes map[string]string) string {
	note, ok := notes[""]
	if ok {
		delete(notes, "")
	}
	for k, v := range notes {
		note += splitFMT + k + splitFMT + v
	}
	return note
}

// update csv data
func (it *Item) update() {
	it.csv.Type = it.Type.String()
	it.csv.Time = it.Time.Format(timeFMT)
	it.csv.Class = it.Class[ClassM]
	if len(it.Class[ClassS]) > 0 {
		it.csv.Class += splitFMT + it.Class[ClassS]
	}
	it.csv.Account = it.Account[AccountM]
	if len(it.Account[AccountS]) > 0 {
		it.csv.Account += splitFMT + it.Account[AccountS]
	}
	it.csv.Amount = it.Amount
	notes := make(map[string]string)
	if len(it.Member) > 0 {
		notes[TagMember.String()] = it.Member
	}
	if len(it.Proj) > 0 {
		notes[TagProj.String()] = it.Proj
	}
	if it.Unit != 0 {
		notes[TagUnit.String()] = fmt.Sprintf("%.2f", it.Unit)
	}
	if it.Deadline.After(it.Time) {
		notes[TagDeadline.String()] = it.Deadline.Format(timeFMT)
	}
	it.csv.Note = it.Note + it.note(notes)
}

func parseItem(data format, year, num int) Item {
	var it Item
	it.csv = data
	it.ID = it.id(year, num)
	it.Type = ParseType(it.csv.Type)
	it.Time = it.csv.time()
	class := it.csv.class()
	if len(class) > ClassS {
		it.Class[ClassS] = class[ClassS]
	}
	if len(class) > ClassM {
		it.Class[ClassM] = class[ClassM]
	}
	it.Amount = it.csv.Amount
	account := it.csv.account()
	if len(account) > AccountS {
		it.Account[AccountS] = account[AccountS]
	}
	if len(account) > AccountM {
		it.Account[AccountM] = account[AccountM]
	}
	notes := it.csv.note()
	for k, v := range notes {
		tag := ParseTag(k)
		switch tag {
		case TagMember:
			it.Member = v
		case TagProj:
			it.Proj = v
		case TagUnit:
			it.Unit, _ = strconv.ParseFloat(v, 64)
		case TagDeadline:
			it.Deadline, _ = time.ParseInLocation(timeFMT, v, time.Local)
		default:
			continue
		}
		delete(notes, k)
	}
	it.Note = it.note(notes)
	return it
}

type items []Item

func (its items) Len() int {
	return len(its)
}

func (its items) Swap(i, j int) {
	its[i], its[j] = its[j], its[i]
}

func (its items) Less(i, j int) bool {
	return its[i].Time.Before(its[j].Time)
}
