package record

import "time"

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
