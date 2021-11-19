package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func LoadConfig(conf *viper.Viper, requiredKeys ...string) error {
	conf.SetConfigName("config")
	conf.SetConfigType("yaml")
	conf.AddConfigPath(".")
	conf.AutomaticEnv()

	if err := conf.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Printf("Config file not found: %v", err)
		default:
			return err
		}
	}

	for _, key := range requiredKeys {
		if !conf.IsSet(key) {
			return fmt.Errorf("required key '%s' is missing in configuration or env variables", key)
		}
	}
	return nil
}
