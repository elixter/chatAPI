package config

import (
	"github.com/spf13/viper"
)

var instance *viper.Viper

func init() {

	vp := viper.New()
	vp.SetConfigFile("config.json")
	vp.SetConfigType("json")

	err := vp.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file not found")
		} else {
			panic(err)
		}
	}
	instance = vp
}

func Config() *viper.Viper {
	return instance
}
