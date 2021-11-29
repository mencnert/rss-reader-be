package cmd

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var (
	webCfg       = &webConfig{}
	dbCfg        = &dbConfig{}
	schedulerCfg = &schedulerConfig{}
)

func loadConfig(vPtr interface{}) error {
	log.Printf("Loading config %T\n", vPtr)
	if rv := reflect.ValueOf(vPtr); rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("config is not pointer: %T", vPtr)
	}
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := bindViperEnvs(vPtr); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Printf("Config file not found: %v\n", err)
		default:
			return err
		}
	}

	if err := viper.Unmarshal(vPtr); err != nil {
		return err
	}

	cfg := reflect.ValueOf(vPtr).Elem().Interface()
	if err := validator.New().Struct(cfg); err != nil {
		return err
	}

	return nil
}

func bindViperEnvs(v interface{}) error {
	cfg := reflect.Indirect(reflect.ValueOf(v)).Interface()
	configType := reflect.TypeOf(cfg)
	configValue := reflect.ValueOf(cfg)

	for i := 0; i < configType.NumField(); i++ {
		fieldType := configType.Field(i)
		fieldValue := configValue.Field(i)
		tagValue, ok := fieldType.Tag.Lookup("mapstructure")
		if !ok {
			return fmt.Errorf("configuration field '%s' is missing required tag 'mapstructure'", fieldType.Name)
		}
		if fieldValue.Kind() == reflect.Struct {
			return fmt.Errorf("err in field '%s': nested structures are not supported in configuration", fieldType.Name)
		}
		if err := viper.BindEnv(tagValue); err != nil {
			return err
		}
	}
	return nil
}

type webConfig struct {
	Port             int      `mapstructure:"PORT"               validate:"required"`
	Username         string   `mapstructure:"LOGIN"              validate:"required"`
	Password         string   `mapstructure:"PASSWORD"           validate:"required"`
	CorsAllowOrigins []string `mapstructure:"CORS_ALLOW_ORIGINS" validate:"min=1,dive,min=1"`
}

type dbConfig struct {
	DBUrl string `mapstructure:"DATABASE_URL" validate:"required,url"`
}

type schedulerConfig struct {
	RssFetchEveryNSecs int `mapstructure:"RSS_FETCH_EVERY_N_SECS" validate:"min=30"`
	CleanDbEveryNHours int `mapstructure:"CLEAN_DB_EVERY_N_HOURS" validate:"min=1"`
}
