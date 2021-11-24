package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("TODO migrate DB", viper.GetString("DATABASE_URL"))
			return nil
		},
	}
	return migrateCmd
}
