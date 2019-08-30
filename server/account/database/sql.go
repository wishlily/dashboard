package database

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type table struct {
	url   string
	name  string
	names []string
	types []string
}

func newTable(url string, name string, st interface{}) (table, error) {
	db := table{url: url, name: name}
	_, err := os.Stat(url)
	if err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(path.Dir(url), 0644)
		}
	}
	if err := db.attributes(reflect.TypeOf(st)); err != nil {
		return db, err
	}
	if err := db.create(); err != nil {
		return db, err
	}
	return db, nil
}

// init table names, types
func (t *table) attributes(typ reflect.Type) error {
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("Can't get %v attributes", typ)
	}
	names := []string{}
	types := []string{}
	for i := 0; i < typ.NumField(); i++ {
		k := typ.Field(i).Tag.Get("db")
		v := typ.Field(i).Tag.Get("type")
		if k == "-" || len(v) == 0 {
			continue
		}
		if len(k) == 0 {
			k = strings.ToLower(typ.Field(i).Name)
		}
		names = append(names, k)
		types = append(types, v)
	}
	t.names = names
	t.types = types
	return nil
}

// create one table in db
func (t table) create() error {
	schema := "("
	for i := 0; i < len(t.names) && i < len(t.types); i++ {
		schema += fmt.Sprintf("%s %s, ", t.names[i], t.types[i])
	}
	schema = strings.TrimRight(schema, ", ")
	schema += ")"
	schema = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", t.name, schema)

	db, err := sqlx.Connect("sqlite3", t.url)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(schema); err != nil {
		return err
	}
	return nil
}

func (t table) insert(data interface{}) error {
	s1, s2 := "(", "("
	for _, v := range t.names {
		s1 += fmt.Sprintf("%s, ", v)
		s2 += fmt.Sprintf(":%s, ", v)
	}
	s1 = strings.TrimRight(s1, ", ")
	s2 = strings.TrimRight(s2, ", ")
	s1 += ")"
	s2 += ")"
	schema := fmt.Sprintf("INSERT INTO %s %s VALUES %s", t.name, s1, s2)
	return t.exec(schema, data)
}

// where: table find cols name
func (t table) update(where string, data interface{}) error {
	schema := ""
	for _, v := range t.names {
		if v == where {
			continue
		}
		schema += fmt.Sprintf("%s = :%s, ", v, v)
	}
	schema = strings.TrimRight(schema, ", ")
	schema = fmt.Sprintf("UPDATE %s SET %s WHERE %s = :%s", t.name, schema, where, where)
	return t.exec(schema, data)
}

func (t table) delete(where string, data interface{}) error {
	schema := fmt.Sprintf("DELETE FROM %s WHERE %s = :%s", t.name, where, where)
	return t.exec(schema, data)
}

func (t table) exec(schema string, data interface{}) error {
	db, err := sqlx.Connect("sqlite3", t.url)
	if err != nil {
		return err
	}
	defer db.Close()

	if res, err := db.NamedExec(schema, data); err != nil {
		return err
	} else if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n <= 0 {
		return fmt.Errorf("database table exec zero: %v", data)
	}
	return nil
}

func (t table) sel(where string, data interface{}) (*sqlx.Rows, error) {
	schema := fmt.Sprintf("SELECT * FROM %s WHERE %s = :%s", t.name, where, where)

	db, err := sqlx.Connect("sqlite3", t.url)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.NamedQuery(schema, data)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (t table) get() (*sqlx.Rows, error) {
	schema := fmt.Sprintf("SELECT * FROM %s", t.name)

	db, err := sqlx.Connect("sqlite3", t.url)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Queryx(schema)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (t table) drop() error {
	schema := fmt.Sprintf("DROP TABLE %s", t.name)

	db, err := sqlx.Connect("sqlite3", t.url)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(schema); err != nil {
		return err
	}
	return nil
}
