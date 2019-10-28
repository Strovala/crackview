package app

import (
	"github.com/spf13/viper"
)

type Language string

type Config struct {
	Languages []Language `json:"languages"`
}

func InitConfig() (*Config, error) {
	config := &Config{
		Languages: viper.GetStringSlice("languages"),
	}
	return config, nil
}
