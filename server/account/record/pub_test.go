package record

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func testCfgPre() {
	viper.SetConfigType("toml")
}

func testCfgClr(t *testing.T) {
	if err := viper.ReadConfig(bytes.NewBuffer([]byte(``))); err != nil { // clear
		t.Fatal(err)
	}
}

func testCfgSet(t *testing.T) {
	if f, err := os.Open("test/csv.toml"); err != nil {
		t.Fatal(err)
	} else {
		defer f.Close()
		if err := viper.ReadConfig(f); err != nil {
			t.Fatal(err)
		}
	}
}

func testCSVSet(t *testing.T, f string, data []byte) {
	if err := ioutil.WriteFile(f, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func testCSVClr(t *testing.T, f string) {
	os.Remove(f)
}
