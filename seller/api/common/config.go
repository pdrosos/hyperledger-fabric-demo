package common

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/logger"
)

func LoadConfig() {
	configName := "config"

	viper.SetConfigName(configName)  // name of config file (without extension)
	viper.AddConfigPath("./config/") // path to look for the config file in
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logger.Log.WithFields(logrus.Fields{
			"configName": configName,
			"error":      err,
		}).Fatal("Unable to load config file")
	}
}
