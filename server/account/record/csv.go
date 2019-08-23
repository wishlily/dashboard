package record

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

var (
	watch string
	cache sync.Map // key ID, value Item
)

// SetPath set csv path : /home/csv/.../
// the path have files like 2015.csv 2016.csv ... year.csv
func SetPath(path string) error {
	path = filepath.Dir(path) // formatting
	watch = path
	if err := watcher(path); err != nil {
		return err
	}
	return nil
}

func watcher(path string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := w.Add(path); err != nil {
		w.Close()
		return err
	}
	go func() {
		log.Infof("Watch path %v start", path)
		defer w.Close()
		for strings.EqualFold(watch, path) {
			select {
			case ev := <-w.Events:
				if ev.Op&fsnotify.Create == fsnotify.Create {
					log.Println("创建文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					log.Println("写入文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("删除文件 : ", ev.Name)
				}
			case err := <-w.Errors:
				log.Errorf("Watch %v error: %v", path, err)
				break
			}
		}
		log.Infof("Watch path %v exit", path)
	}()
	return nil
}

type file struct {
	start int64
	end   int64
	name  string
}

// SetRange set load csv item data time range
func SetRange(start, end int64) error {
	// TODO:
	return nil
}

// type _Bills []Bill

// func (m _Bills) Len() int {
// 	return len(m)
// }

// func (m _Bills) Swap(i, j int) {
// 	m[i], m[j] = m[j], m[i]
// }

// func (m _Bills) Less(i, j int) bool {
// 	return m[i].Time < m[j].Time
// }

// type CSV struct {
// 	file  string
// 	cache map[string]Bill
// 	data  _Bills
// }

// func NewCSV(f string) *CSV {
// 	var c CSV
// 	c.file = f
// 	_, err := os.Stat(f)
// 	if err != nil {
// 		if !os.IsExist(err) {
// 			os.MkdirAll(path.Dir(f), 0644)
// 			os.Create(f)
// 			c.writeHeader()
// 		}
// 	}
// 	return &c
// }

// func (c *CSV) Data() ([]Bill, error) {
// 	if err := c.readAll(); err != nil {
// 		return nil, err
// 	}
// 	return c.data, nil
// }

// func (c *CSV) Get(id string) (Bill, bool) {
// 	v, ok := c.cache[id]
// 	return v, ok
// }

// func (c *CSV) Add(data Bill) error {
// 	return c.write(data)
// }

// func (c *CSV) Del(id string) error {
// 	if _, ok := c.cache[id]; ok {
// 		delete(c.cache, id)
// 		c.updateData()
// 		return c.writeAll()
// 	}
// 	return errcode.NFIND
// }

// func (c *CSV) Change(data Bill) error {
// 	if _, ok := c.cache[data.ID]; ok {
// 		c.cache[data.ID] = data
// 		c.updateData()
// 		return c.writeAll()
// 	}
// 	return errcode.NFIND
// }

// func (c *CSV) genID(n int) string {
// 	filenameall := path.Base(c.file)
// 	suffix := path.Ext(filenameall)
// 	filename := strings.TrimSuffix(filenameall, suffix)
// 	return fmt.Sprintf("%v%v%d", filename, time.Now().Unix(), n)
// }

// func (c *CSV) updateData() {
// 	var data _Bills
// 	for _, v := range c.cache {
// 		data = append(data, v)
// 	}
// 	sort.Sort(data)
// 	c.data = data
// }

// func (c *CSV) update(data []Bill) {
// 	m := make(map[string]Bill)
// 	for i, v := range data {
// 		v.ID = c.genID(i)
// 		m[v.ID] = v
// 	}
// 	c.cache = m
// 	c.updateData()
// }

// func (c *CSV) readAll() error {
// 	bytes, err := ioutil.ReadFile(c.file)
// 	if err != nil {
// 		return err
// 	}
// 	var data []Bill
// 	if err := csvutil.Unmarshal(bytes, &data); err != nil {
// 		return err
// 	}
// 	// generate ID & Map Data
// 	c.update(data)
// 	return nil
// }

// func (c *CSV) writeAll() error {
// 	content, err := csvutil.Marshal(c.data)
// 	if err != nil {
// 		return err
// 	}
// 	if err := ioutil.WriteFile(c.file, content, 0644); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *CSV) write(v Bill) error {
// 	f, err := os.OpenFile(c.file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	w := csv.NewWriter(f)
// 	enc := csvutil.NewEncoder(w)
// 	enc.AutoHeader = false
// 	if err := enc.Encode(v); err != nil {
// 		return err
// 	}
// 	w.Flush()
// 	if err := w.Error(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *CSV) writeHeader() error {
// 	header, err := csvutil.Header(Bill{}, "csv")
// 	if err != nil {
// 		return err
// 	}
// 	f, err := os.OpenFile(c.file, os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	fmt.Fprintf(f, "%s\n", strings.Join(header, ","))
// 	return nil
// }
