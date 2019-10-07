package config

import (
	"errors"
	"github.com/chromz/replicator/pkg/log"
	"github.com/spf13/viper"
)

type Config struct {
	Options Options
	Server  Server
}

type Options struct {
	Directory   string
	Module      string
	PollingRate int `mapstructure:"polling-rate"`
}

type Server struct {
	Name    string
	Port    string
	Address string
}

var loadedConfig Config

func LoadConfig(name string) error {
	viper.SetConfigType("toml")
	if name != "" {
		viper.SetConfigFile(name)
	} else {
		viper.AddConfigPath("/etc/replicator/")
		viper.AddConfigPath("./configs")
		viper.SetConfigName("config")
	}
	viper.SetDefault("port", 4690)
	viper.SetDefault("directory", "/var/replicator/sync/")
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

func Port() string {
	return loadedConfig.Server.Port
}

func Directory() string {
	return loadedConfig.Options.Directory
}

func Module() string {
	return loadedConfig.Options.Module
}

func RsyncServer() Server {
	return loadedConfig.Server
}

func Host() string {
	return loadedConfig.Server.Address + ":" + loadedConfig.Server.Port
}
