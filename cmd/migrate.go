package cmd

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/spf13/cobra"
)

var (
	migrateDb *migrate.Migrate
)

func newMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use: "migrate",
	}
	migrateCmd.AddCommand(newMigrateUpCmd())
	migrateCmd.AddCommand(newMigrateDownCmd())
	migrateCmd.AddCommand(newMigrateToCmd())

	return migrateCmd
}

func newMigrateUpCmd() *cobra.Command {
	migrateUpCmd := &cobra.Command{
		Use:     "up",
		PreRunE: prepareMigration,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("DB migration: start")
			defer log.Println("DB migration: done")
			err := migrateDb.Up()
			if err == migrate.ErrNoChange {
				log.Println("DB migration: no change")
				return nil
			}
			return err
		},
	}

	return migrateUpCmd
}

func newMigrateDownCmd() *cobra.Command {
	migrateDownCmd := &cobra.Command{
		Use:     "down",
		PreRunE: prepareMigration,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("DB migration: start")
			defer log.Println("DB migration: done")
			err := migrateDb.Down()
			if err == migrate.ErrNoChange {
				log.Println("DB migration: no change")
				return nil
			}
			return err
		},
	}

	return migrateDownCmd
}

func newMigrateToCmd() *cobra.Command {
	var dbVersion uint
	migrateToCmd := &cobra.Command{
		Use:     "to",
		PreRunE: prepareMigration,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("DB migration: start")
			defer log.Println("DB migration: done")
			err := migrateDb.Migrate(dbVersion)
			if err == migrate.ErrNoChange {
				log.Println("DB migration: no change")
				return nil
			}
			if os.IsNotExist(err) {
				log.Printf("DB migration: err unable to find file for specified version: %d\n", dbVersion)
				return err
			}
			return err
		},
	}
	migrateToCmd.Flags().UintVarP(&dbVersion, "version", "v", 0, "DB migration version")
	migrateToCmd.MarkFlagRequired("version")

	return migrateToCmd
}

func prepareMigration(cmd *cobra.Command, args []string) error {
	log.Println("DB migration: initialize migration")
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)

	if err != nil {
		return err
	}

	migrateDb = m
	return nil
}
