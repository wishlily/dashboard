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
	ID       string
	Type     Type
	Time     time.Time
	Class    [classN]string
	Amount   float64
	Account  [accountN]string
	Member   string
	Proj     string // project & id
	Unit     int64
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
		note += "#" + k + "#" + v
	}
	return note
}

func parseItem(data format, year, num int) Item {
	var it Item
	it.csv = data
	it.ID = it.id(year, num)
	it.Type = ParseType(it.csv.Type)
	it.Time = it.csv.time()
	class := it.csv.class()
	if len(class) >= ClassS {
		it.Class[ClassS] = class[ClassS]
	}
	if len(class) >= ClassM {
		it.Class[ClassM] = class[ClassM]
	}
	it.Amount = it.csv.Amount
	account := it.csv.account()
	if len(account) >= AccountS {
		it.Account[AccountS] = account[AccountS]
	}
	if len(account) >= AccountM {
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
			it.Unit, _ = strconv.ParseInt(v, 10, 64)
		case TagDeadline:
			it.Deadline, _ = time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		default:
			continue
		}
		delete(notes, k)
	}
	it.Note = it.note(notes)
	return it
}
