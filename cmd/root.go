package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCmd(requiredKeys []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "rss-reader",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setupConfig(requiredKeys); err != nil {
				return err
			}
			return nil
		},
	}
	rootCmd.AddCommand(newWebCmd())
	rootCmd.AddCommand(newMigrateCmd())

	return rootCmd
}

func Execute(requiredKeys ...string) error {
	rootCmd := newRootCmd(requiredKeys)
	return rootCmd.Execute()
}

func setupConfig(requiredKeys []string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Printf("Config file not found: %v", err)
		default:
			return err
		}
	}

	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			return fmt.Errorf("required key '%s' is missing in configuration or env variables", key)
		}
	}
	return nil
}
