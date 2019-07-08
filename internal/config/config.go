package config

import (
	"github.com/chromz/replicator/pkg/log"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/replicator/")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		log.ErrorLog("Error reading config file", err)
	}
}
