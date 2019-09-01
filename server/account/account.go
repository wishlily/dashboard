package account

import (
	"github.com/spf13/viper"
	db "github.com/wishlily/dashboard/server/account/database"
	rd "github.com/wishlily/dashboard/server/account/record"
)

// Init account config by viper
func Init() error {
	path := viper.GetString("csv.path")
	if len(path) == 0 {
		path = "db"
	}
	if err := rd.SetPath(path); err != nil {
		return err
	}
	url := viper.GetString("db.url")
	if len(url) == 0 {
		url = "db/database.db"
	}
	if err := db.SetURL(url); err != nil {
		return err
	}
	return nil
}

// AddRecord add one record
func AddRecord(item rd.Item) error {
	record := record{item}
	switch item.Type {
	case rd.TypeI:
		return record.io(rd.Add, record.add)
	case rd.TypeO:
		return record.io(rd.Add, record.sub)
	case rd.TypeR:
		return record.r(rd.Add, record.sub, record.add)
	}
	return nil
}

type record struct {
	rd.Item
}

func (r record) add(acct *db.Account) {
	acct.Unit += r.Unit
	acct.Input += r.Amount
}

func (r record) sub(acct *db.Account) {
	acct.Unit -= r.Unit
	acct.Input -= r.Amount
}

// f : rd.Add rd.Del rd.Chg
func (r record) io(f func(rd.Item) error, opt func(acct *db.Account)) error {
	id := r.Account[rd.AccountM]
	acct, err := db.GetAccount().Sel(id)
	if err != nil {
		return err
	}
	opt(&acct)
	if err := f(r.Item); err != nil {
		return err
	}
	return db.GetAccount().Chg(acct)
}

func (r record) r(f func(rd.Item) error, optM func(acct *db.Account), optS func(acct *db.Account)) error {
	idM := r.Account[rd.AccountM]
	acctM, err := db.GetAccount().Sel(idM)
	if err != nil {
		return err
	}
	idS := r.Account[rd.AccountS]
	acctS, err := db.GetAccount().Sel(idS)
	optM(&acctM)
	optS(&acctS)
	if err := f(r.Item); err != nil {
		return err
	}
	if err := db.GetAccount().Chg(acctM); err != nil {
		return err
	}
	return db.GetAccount().Chg(acctS)
}
