package api

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port int `json:"port"`
}

func InitConfig() (*Config, error) {
	config := &Config{
		Port: viper.GetInt("port"),
	}
	return config, nil
}
