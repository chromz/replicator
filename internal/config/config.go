package config

import (
	"errors"
	"github.com/chromz/replicator/pkg/log"
	"github.com/spf13/viper"
)

// Config struct that represents the config file
type Config struct {
	Options Options
	Server  Server
}

// Options is the "option" section of the config file
type Options struct {
	Directory   string
	Module      string
	PollingRate int `mapstructure:"polling-rate"`
}

// Server is the server section of the config file
type Server struct {
	Name    string
	Address string
}

var loadedConfig Config

// LoadConfig loads the config file into memory
func LoadConfig() error {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/replicator/")
	viper.AddConfigPath("./configs")
	viper.SetDefault("port", 4690)
	viper.SetDefault("directory", "/var/replicator/sync/")
	viper.SetDefault("polling-rate", 5)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Error reading config file", err)
		return err
	}
	if err := viper.Unmarshal(&loadedConfig); err != nil {
		log.Error("Invalid config file", err)
		return err
	}
	if loadedConfig.Options.Module == "" {
		err := errors.New("no rsync module")
		log.Error("Invalid config file", err)
		return err
	}
	if loadedConfig.Server.Address == "" {
		err := errors.New("no server address")
		log.Error("Invalid config file", err)
		return err

	}
	dir := loadedConfig.Options.Directory
	if dir == "" {
		err := errors.New("no server address")
		log.Error("Invalid config file", err)
		return err
	}
	if dir[len(dir)-1:] != "/" {
		loadedConfig.Options.Directory += "/"
	}
	log.Info("Config file loaded")
	return nil
}

// Directory returns the directory to synchronize
func Directory() string {
	return loadedConfig.Options.Directory
}

// Module gets the rsync module
func Module() string {
	return loadedConfig.Options.Module
}

// RsyncServer gets the rsync server address
func RsyncServer() Server {
	return loadedConfig.Server
}

// PollingRate gets the polling rate to update
func PollingRate() int {
	return loadedConfig.Options.PollingRate
}
