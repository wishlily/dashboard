package config

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestInit(t *testing.T) {
	// Overrides
	if v := viper.GetString("Title"); !strings.EqualFold(v, "TOML Example") {
		t.Fatalf("%v != %v", v, "TOML Example")
	}
	// nested keys
	if v := viper.GetString("servers.alpha.ip"); !strings.EqualFold(v, "10.0.0.1") {
		t.Fatalf("%v != %v", v, "10.0.0.1")
	}
	// get null
	if v := viper.Get("NIL"); v != nil {
		t.Fatalf("%v != nil", v)
	}
}

func TestParse(t *testing.T) {
	Parse()

	if v := viper.GetString("log"); !strings.EqualFold(v, "info") {
		t.Fatalf("%v != %v", v, "info")
	}
	if v := viper.GetBool("database.enabled"); !v {
		t.Fatalf("Lost init config")
	}
}

func TestRead(t *testing.T) {
	const (
		file = "nil.toml"
	)

	if err := Read(file); err == nil {
		t.Fatalf("Read should be error")
	}

	f, _ := os.Create(file)
	f.Close()
	defer os.Remove(file)

	if err := Read(file); err != nil {
		t.Fatalf("Read error %v", err)
	}
	if v := viper.GetBool("database.enabled"); v {
		t.Fatalf("Data should be false")
	}
	if err := Read("config.toml"); err != nil {
		t.Fatalf("Read error %v", err)
	}
	if v := viper.GetBool("database.enabled"); !v {
		t.Fatalf("Data should be true")
	}
}
