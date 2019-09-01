package database

const (
	nameAccount = "account"
	nameBorrow  = "borrow"
	nameLend    = "lend"
)

var (
	cacheA AccountTable
	cacheB DebitTable
	cacheL DebitTable
)

// SetURL set database url
func SetURL(url string) error {
	var err error
	cacheA, err = newAccountTable(url, nameAccount)
	if err != nil {
		return err
	}
	cacheB, err = newDebitTable(url, nameBorrow)
	if err != nil {
		return err
	}
	cacheL, err = newDebitTable(url, nameLend)
	if err != nil {
		return err
	}
	return nil
}

// GetAccount get db table account
func GetAccount() AccountTable {
	return cacheA
}

// GetBorrow get db table borrow debit
func GetBorrow() DebitTable {
	return cacheB
}

// GetLend get db table lend debit
func GetLend() DebitTable {
	return cacheL
}
