package account

import (
	"fmt"

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
func AddRecord(data rd.Item) error {
	record := record{data}
	switch data.Type {
	case rd.TypeI:
		return record.io(rd.Add, record.add)
	case rd.TypeO:
		return record.io(rd.Add, record.sub)
	case rd.TypeR:
		return record.r(rd.Add, record.sub, record.add)
	case rd.TypeB:
		return record.debit(record.add, db.GetBorrow())
	case rd.TypeL:
		return record.debit(record.sub, db.GetLend())
	}
	return fmt.Errorf("Can't support %v type record add", data.Type)
}

// DelRecord del one record
func DelRecord(data rd.Item) error {
	record := record{data}
	switch data.Type {
	case rd.TypeI:
		return record.io(rd.Del, record.sub)
	case rd.TypeO:
		return record.io(rd.Del, record.add)
	case rd.TypeR:
		return record.r(rd.Del, record.add, record.sub)
	}
	return fmt.Errorf("Can't support %v type record add", data.Type)
}

// ChgRecord change one record
func ChgRecord(data rd.Item) error {
	old, err := rd.Sel(data)
	if err != nil {
		return err
	}
	if err := DelRecord(old); err != nil {
		return err
	}
	return AddRecord(data)
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
	if acct.Unit < 0 {
		acct.Unit = 0 // Unit not zero
	}
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

func (r record) debit(opt func(acct *db.Account), tab db.DebitTable) error {
	member := r.Member
	if len(member) == 0 {
		return fmt.Errorf("Should be have member")
	}
	id := r.Account[rd.AccountM]
	acct, err := db.GetAccount().Sel(id)
	if err != nil {
		return err
	}
	opt(&acct)
	data := db.Debit{Name: member, Amount: r.Amount, Note: r.Note, Deadline: r.Deadline}
	if err := tab.Add(data); err != nil {
		return err
	}
	if err := rd.Add(r.Item); err != nil {
		return err
	}
	return db.GetAccount().Chg(acct)
}
