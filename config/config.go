package config

import (
	"chatting/logger"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

package config

import (
	"ImageRemover/logging"
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
			logger.Log.Errorf("Config file not found - file name : [%s]", vp.ConfigFileUsed())
		} else {
			log.Error(err)
		}
	}
	instance = vp
}

func Config() *viper.Viper {
	return instance
}