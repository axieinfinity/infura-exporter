package config

import (
	"fmt"
	"github.com/namcxn/infura-exporter/logger"
	"github.com/spf13/viper"
	"strings"
)

type Schema struct {
	Infura struct {
		Url   string `mapstructure:"url"`
		Token string `mapstructure:"token"`
	} `mapstructure:"infura"`
}

var (
	logs   = logger.GetLogger("Config")
	Config Schema
)

func init() {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")           // Look for config in current directory.
	config.AddConfigPath("pkg/config/") // Look for config needed for tests.
	config.AddConfigPath("config/")     // Optionally look for config in the working directory.
	config.AddConfigPath("../config/")
	config.AddConfigPath("../../")

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = config.Unmarshal(&Config)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	logs.Infof("Config: %+v", Config)
}
