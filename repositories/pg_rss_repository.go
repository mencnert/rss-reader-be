package repositories

import (
	"database/sql"
	"log"
	"rss-reader/rss"

	_ "github.com/lib/pq"
)

type pgRssRepository struct {
	url string
	db  *sql.DB
}

func (pg *pgRssRepository) Open() error {
	log.Println("Connecting to DB")
	db, err := sql.Open("postgres", pg.url)
	if err != nil {
		return err
	}
	pg.db = db
	log.Println("Test DB connection")
	if err := pg.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (pg *pgRssRepository) Close() error {
	log.Println("Closing DB connection")
	return pg.db.Close()
}

func (pg *pgRssRepository) GetAll() ([]RssDTO, error) {
	sqlQuery := `SELECT rss_id, url, rank, title, viewed, saved FROM rss`

	rows, err := pg.db.Query(sqlQuery)
	if err != nil {
		return []RssDTO{}, err
	}
	defer rows.Close()
	result := make([]RssDTO, 0)
	for rows.Next() {
		var dto RssDTO
		if err := rows.Scan(&dto.Id, &dto.Url, &dto.Rank, &dto.Title, &dto.Viewed, &dto.Saved); err != nil {
			return []RssDTO{}, err
		}
		result = append(result, dto)
	}
	return result, nil
}

func (pg *pgRssRepository) SaveOrUpdateAll(rssEntries []rss.RssEntry) error {
	sqlQuery := `
		INSERT INTO rss (url, rank, title, last_fetch)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (url) DO UPDATE
		  SET rank = excluded.rank,
		  title = excluded.title,
		  last_fetch = NOW();`

	stmt, err := pg.db.Prepare(sqlQuery)
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

func (pg *pgRssRepository) Update(rssDto RssDTO) error {
	sqlQuery := `
		UPDATE rss SET
			viewed = $2,
			saved  = $3
		WHERE rss_id = $1`

	stmt, err := pg.db.Prepare(sqlQuery)
	if err != nil {
		log.Printf("Error during preparation of query: %v\n", err)
		return err
	}
	if _, err := stmt.Exec(rssDto.Id, rssDto.Viewed, rssDto.Saved); err != nil {
		log.Printf("Error during execution of query: %v\n", err)
		return err
	}

	return nil
}

func NewPgRssRepository(url string) *pgRssRepository {
	return &pgRssRepository{url: url}
}
