package database

var (
	cache database
)

type database struct {
	account accountTable
	borrow  debitTable
	lend    debitTable
}

func init() {
	const defaultURL = "database.db"
	cache.account, _ = newAccountTable(defaultURL)
	cache.borrow, _ = newBorrowTable(defaultURL)
	cache.lend, _ = newLendTable(defaultURL)
}

// SetURL set database url
func SetURL(url string) error {
	account, err := newAccountTable(url)
	if err != nil {
		return err
	}
	cache.account = account
	obj, err := newBorrowTable(url)
	if err != nil {
		return err
	}
	cache.borrow = obj
	obj, err = newLendTable(url)
	if err != nil {
		return err
	}
	cache.lend = obj
	return nil
}
