package record

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

var (
	cache database
)

type database struct {
	path string
	buf  sync.Map // key ID, value Item
	min  int      // load files year range min
	max  int      // load files year range max
}

// SetPath set csv path : /home/csv/.../
// the path have files like 2015.csv 2016.csv ... year.csv
func SetPath(path string) error {
	return cache.setPath(path)
}

// Get [start, end] time records from csv files
// start|end: "yyyy-mm-dd hh:mm:ss"
func Get(start, end string) ([]Item, error) {
	t1, err := time.ParseInLocation(timeFMT, start, time.Local)
	if err != nil {
		return nil, err
	}
	t2, err := time.ParseInLocation(timeFMT, end, time.Local)
	if err != nil {
		return nil, err
	}
	if err := cache.setRange(t1, t2); err != nil {
		return nil, err
	}
	var data items
	cache.buf.Range(func(k interface{}, v interface{}) bool {
		item := v.(Item) // no check data type
		if (item.Time.After(t1) || item.Time.Equal(t1)) &&
			(item.Time.Before(t2) || item.Time.Equal(t2)) {
			data = append(data, item) // check time
		}
		return true
	})
	sort.Sort(data)
	return data, nil
}

// Add one item in csv files
func Add(data Item) error {
	year := data.Time.Year()
	if year <= 2000 {
		return fmt.Errorf("Add item time year can't less 2000 : %v", year)
	}
	w := writer{cache.csv(year)}
	data.update()
	return w.append(data.csv)
}

// Del the item in csv files
func Del(data Item) error {
	if _, ok := cache.buf.Load(data.ID); !ok {
		return fmt.Errorf("Del item not found ID %v", data.ID)
	}
	cache.buf.Delete(data.ID)
	return cache.save(data.Time.Year())
}

// Chg (change) the item in csv files
// return old item data
func Chg(data Item) (old Item, err error) {
	it, ok := cache.buf.Load(data.ID)
	if !ok {
		return old, fmt.Errorf("Chg not found ID %v", data.ID)
	}
	old = it.(Item) // no check data type
	if err := Del(old); err != nil {
		return old, err
	}
	return old, Add(data)
}

func (db *database) setPath(path string) error {
	path, err := filepath.Abs(path) // formatting
	if err != nil {
		return err
	}
	db.path = path
	return nil
}

func (db *database) csv(year int) string {
	if len(db.path) == 0 {
		db.path = "."
	}
	return fmt.Sprintf("%s/%04d.csv", db.path, year)
}

func (db *database) setRange(start, end time.Time) error {
	if end.Before(start) {
		return fmt.Errorf("SetRange can't set end < start ts: %v,%v", end, start)
	}
	db.min = start.Year()
	db.max = end.Year()
	return db.load()
}

// load files data in buf
func (db *database) load() error {
	if db.min == 0 || db.max == 0 {
		return fmt.Errorf("Load min %v max %v", db.min, db.max)
	}
	for year := db.min; year <= db.max; year++ {
		r := reader{db.csv(year)}
		data, err := r.all()
		if err != nil {
			return err
		}
		for i, v := range data {
			it := parseItem(v, year, i)
			db.buf.Store(it.ID, it) // store buf
		}
	}
	return nil
}

// save files data in buf
func (db *database) save(year int) error {
	f := db.csv(year)
	w := writer{f}

	var data items
	cache.buf.Range(func(k interface{}, v interface{}) bool {
		item := v.(Item)              // no check data type
		if item.Time.Year() == year { // check time
			data = append(data, item)
		}
		return true
	})
	sort.Sort(data)

	var store []format
	for _, v := range data {
		store = append(store, v.csv)
	}
	return w.all(store)
}
