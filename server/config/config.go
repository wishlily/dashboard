package config

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("title", "dashboard")

	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.dashboard")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("Config file not found")
		} else {
			log.Panicf("Fatal error config file: %s\n", err)
		}
	} else {
		log.Info("Load config file success")
	}
}

// Parse : bind flags in config
func Parse() {
	flag.String("log", "info", "Set log level - Trace, Debug, Info, Warning, Error, Fatal and Panic")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

// Read path config
func Read(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return viper.ReadConfig(f)
}
