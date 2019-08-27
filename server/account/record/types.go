package record

import (
	"strings"

	"github.com/spf13/viper"
)

// Type Item.Type data type
type Type int

const (
	// TypeI : income
	TypeI Type = iota
	// TypeO : pay out
	TypeO
	// TypeR : transfer
	TypeR
	// TypeX : update
	TypeX
	// TypeB : borrow
	TypeB
	// TypeL : lend
	TypeL
	typeEnd
)

func (t Type) key(section, key string) string {
	url := section + "." + key
	v := viper.GetString(url)
	if viper.IsSet(url) && len(v) != 0 {
		return v
	}
	return key
}

func (t Type) String() string {
	const section = "csv.types"
	switch t {
	case TypeI:
		return t.key(section, "I")
	case TypeO:
		return t.key(section, "O")
	case TypeR:
		return t.key(section, "R")
	case TypeX:
		return t.key(section, "X")
	case TypeB:
		return t.key(section, "B")
	case TypeL:
		return t.key(section, "L")
	}
	return t.key(section, "U") // unkownn
}

// ParseType : Parse string to Type
func ParseType(v string) Type {
	var t Type
	for t = TypeI; t < typeEnd; t++ {
		if strings.EqualFold(v, t.String()) {
			return t
		}
	}
	return typeEnd
}
