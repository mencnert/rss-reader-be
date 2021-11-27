package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"rss-reader/rss"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rssRepo *PgRssRepository
	Config  = &Configuration{}
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
			rssRepo = NewPgRssRepository(db)
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			closeConn := func(c DbCloser) error {
				log.Println("Closing DB connection")
				return c.Close()
			}

			return closeConn(rssRepo)
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

// ############# repository

type DbCloser interface {
	Close() error
}

type RssRepository interface {
	SaveOrUpdateAll(rssEntries []rss.RssEntry) error
	GetAll() ([]RssDTO, error)
	UpdateViewedById(id int, viewed bool) error
}

type PgRssRepository struct {
	DB *sql.DB
}

type RssDTO struct {
	Id     int    `json:"id"`
	Url    string `json:"url"`
	Rank   int    `json:"rank"`
	Title  string `json:"title"`
	Viewed bool   `json:"viewed"`
}

func (pg PgRssRepository) Close() error {
	return pg.DB.Close()
}

func (pg PgRssRepository) GetAll() ([]RssDTO, error) {
	sqlQuery := `SELECT rss_id, url, rank, title, viewed FROM rss`

	rows, err := pg.DB.Query(sqlQuery)
	if err != nil {
		return []RssDTO{}, err
	}
	defer rows.Close()
	result := make([]RssDTO, 0)
	for rows.Next() {
		var dto RssDTO
		if err := rows.Scan(&dto.Id, &dto.Url, &dto.Rank, &dto.Title, &dto.Viewed); err != nil {
			return []RssDTO{}, err
		}
		result = append(result, dto)
	}
	return result, nil
}

func (pg PgRssRepository) SaveOrUpdateAll(rssEntries []rss.RssEntry) error {
	sqlQuery := `
		INSERT INTO rss (url, rank, title, last_fetch)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (url) DO UPDATE
		  SET rank = excluded.rank,
		  title = excluded.title,
		  last_fetch = NOW();`

	stmt, err := pg.DB.Prepare(sqlQuery)
	if err != nil {
		log.Printf("Error during preparation of query: %v\n", err)
		return err
	}
	for _, entry := range rssEntries {
		if _, err := stmt.Exec(entry.Url, entry.Rank, entry.Title); err != nil {
			log.Printf("Error during execution of query: %v\n", err)
			return err
		}
	}

	return nil
}

func (pg PgRssRepository) UpdateViewedById(id int, viewed bool) error {
	sqlQuery := `
		UPDATE rss SET
			viewed = $1
		WHERE rss_id = $2`

	stmt, err := pg.DB.Prepare(sqlQuery)
	if err != nil {
		log.Printf("Error during preparation of query: %v\n", err)
		return err
	}
	if _, err := stmt.Exec(viewed, id); err != nil {
		log.Printf("Error during execution of query: %v\n", err)
		return err
	}

	return nil
}

func NewPgRssRepository(db *sql.DB) *PgRssRepository {
	return &PgRssRepository{db}
}
