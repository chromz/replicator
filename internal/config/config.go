package config

import (
	"github.com/chromz/replicator/pkg/log"
	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/replicator/")
	viper.AddConfigPath("./configs")
	viper.SetDefault("port", 4690)
	viper.SetDefault("directory", "./")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Error reading config file", err)
		return err
	}
	log.Info("Config file loaded")
	return nil
}

func GetDirectory() string {
	return viper.GetString("directory")
}
