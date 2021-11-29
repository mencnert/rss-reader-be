package cmd

import (
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "rss-reader",
	}
	rootCmd.AddCommand(newWebCmd())
	rootCmd.AddCommand(newMigrateCmd())

	return rootCmd
}

func Execute() error {
	rootCmd := newRootCmd()
	return rootCmd.Execute()
}
