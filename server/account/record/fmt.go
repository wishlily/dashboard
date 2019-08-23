package record

import (
	"strings"
)

type format struct {
	Type    string  `csv:"type"`
	Time    string  `csv:"time"`    // 2006-01-02 15:04:05
	Class   string  `csv:"class"`   // # split
	Account string  `csv:"account"` // # split
	Amount  float64 `csv:"amount"`
	Note    string  `csv:"note"` // ## tags
}

// #
func (f format) split(data string) []string {
	ss := strings.SplitN(data, "#", -1)
	isNull := func(r rune) bool {
		return r == '\n' || r == ' ' || r == '\t'
	}
	for i, v := range ss {
		ss[i] = strings.TrimFunc(v, isNull)
	}
	return ss
}

// xxx#tag1#data1#tag2#data2
func (f format) note() map[string]string {
	ss := f.split(f.Note)
	note := make(map[string]string)
	note[""] = ss[0] // xxx
	for i := 1; i < len(ss)-1; i += 2 {
		if len(ss[i]) > 0 {
			note[ss[i]] = ss[i+1]
		}
	}
	return note
}

func (f format) class() []string {
	return f.split(f.Class)
}

func (f format) account() []string {
	return f.split(f.Account)
}
