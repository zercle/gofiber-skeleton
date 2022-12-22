package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(configName string) (err error) {
	if len(configName) == 0 {
		configName = "dev"
	}
	viper.SetConfigName(configName)  // name of config file (without extension)
	viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // path to look for the config file in
	viper.AddConfigPath(".")         // optionally look for config in the working directory
	err = viper.ReadInConfig()       // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return
}
