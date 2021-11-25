package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	DB     *sql.DB
	Config = &Configuration{}
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "rss-reader",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setupConfig(); err != nil {
				return err
			}
			log.Println("Connecting to DB")
			db, err := sql.Open("postgres", Config.DBUrl)
			if err != nil {
				return err
			}
			DB = db
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Closing DB connection")
			return DB.Close()
		},
	}
	rootCmd.AddCommand(newWebCmd())
	rootCmd.AddCommand(newMigrateCmd())

	return rootCmd
}

func Execute() error {
	rootCmd := newRootCmd()
	return rootCmd.Execute()
}

func setupConfig() error {
	log.Println("Preparing viper config")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := bindViperEnvs(); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Printf("Config file not found: %v", err)
		default:
			return err
		}
	}

	if err := viper.Unmarshal(Config); err != nil {
		return err
	}

	v := validator.New()
	if err := v.Struct(Config); err != nil {
		return err
	}

	return nil
}

func bindViperEnvs() error {
	cfg := Configuration{}
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

type Configuration struct {
	Port               int      `mapstructure:"PORT"                   validate:"required"`
	Username           string   `mapstructure:"LOGIN"                  validate:"required"`
	Password           string   `mapstructure:"PASSWORD"               validate:"required"`
	RssFetchEveryNSecs int      `mapstructure:"RSS_FETCH_EVERY_N_SECS" validate:"min=30"`
	CleanDbEveryNHours int      `mapstructure:"CLEAN_DB_EVERY_N_HOURS" validate:"min=1"`
	DBUrl              string   `mapstructure:"DATABASE_URL"           validate:"required,url"`
	CorsAllowOrigins   []string `mapstructure:"CORS_ALLOW_ORIGINS"     validate:"min=1,dive,min=1"`
}
