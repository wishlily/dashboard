package record

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
)

const (
	timeFMT  = "2006-01-02 15:04:05"
	splitFMT = "#"
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
	ss := strings.SplitN(data, splitFMT, -1)
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

func (f format) time() time.Time {
	t, _ := time.ParseInLocation(timeFMT, f.Time, time.Local)
	return t
}

func (f format) hash() string {
	ss := fmt.Sprintf("%v", f)
	sum := sha1.Sum([]byte(ss))
	return fmt.Sprintf("%x", sum)
}

type reader struct {
	f string
}

// if file not exist create one & header
func (r reader) is() error {
	if _, err := os.Stat(r.f); err != nil && !os.IsExist(err) {
		w := writer{r.f}
		os.MkdirAll(path.Dir(r.f), 0644)
		os.Create(r.f)
		return w.header()
	}
	return nil
}

func (r reader) all() ([]format, error) {
	if err := r.is(); err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadFile(r.f)
	if err != nil {
		return nil, err
	}
	var data []format
	if err := csvutil.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}

type writer struct {
	f string
}

func (w writer) header() error {
	hr, err := csvutil.Header(format{}, "csv")
	if err != nil {
		return err
	}
	f, err := os.OpenFile(w.f, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "%s\n", strings.Join(hr, ","))
	return nil
}

func (w writer) append(data format) error {
	r := reader{w.f}
	if err := r.is(); err != nil {
		return err
	}
	f, err := os.OpenFile(w.f, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	ptr := csv.NewWriter(f)
	enc := csvutil.NewEncoder(ptr)
	enc.AutoHeader = false
	if err := enc.Encode(data); err != nil {
		return err
	}
	ptr.Flush()
	if err := ptr.Error(); err != nil {
		return err
	}
	return nil
}

func (w writer) all(data []format) error {
	content, err := csvutil.Marshal(data)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(w.f, content, 0644); err != nil {
		return err
	}
	return nil
}
