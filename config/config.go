package config

import (
	"managedata/model"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var AppConfig *model.Configurations

func GetConfigurations() {
	var configPath string

	// Get current directory path
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// Define config file path
	configPath = basepath + "/data/"

	// Set Viper configurations
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Config file error: ", err)
	}

	// Unmarshal configurations into struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatal("Error unmarshalling config JSON: ", err)
	}
}
