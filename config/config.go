package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func New() Manager {
	viper.AddConfigPath(path)
	viper.SetConfigName(local)
	viper.SetConfigType(yaml)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("error while reading config file: %s", err))
	}

	var config config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("error while unmarshalling config file: %s", err))
	}

	global = &manager{config: &config}
	return global
}

const (
	path  = "./basket-service/.config"
	local = "local"
	yaml  = "yaml"
)

func Global() Manager {
	if global == nil {
		return New()
	}

	return global
}
