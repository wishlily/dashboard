package database

import (
	"os"
	"testing"
)

func TestNull(t *testing.T) {
	if _, err := GetAccount().Get(); err == nil {
		t.Fatalf("Should be error")
	}
	if _, err := GetBorrow().Get(); err == nil {
		t.Fatalf("Should be error")
	}
	if _, err := GetLend().Get(); err == nil {
		t.Fatalf("Should be error")
	}
}

func TestSetURL(t *testing.T) {
	const URL = "test.db"
	defer os.Remove(URL)

	if err := SetURL(URL); err != nil {
		t.Fatal(err)
	}
	if _, err := GetAccount().Get(); err != nil {
		t.Fatal(err)
	}
	if _, err := GetBorrow().Get(); err != nil {
		t.Fatal(err)
	}
	if _, err := GetLend().Get(); err != nil {
		t.Fatal(err)
	}
}
