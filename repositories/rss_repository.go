package repositories

import (
	"rss-reader/rss"
)

type RssRepository interface {
	Open() error
	SaveOrUpdateAll(rssEntries []rss.RssEntry) error
	GetAll() ([]RssDTO, error)
	Update(RssDTO) error
	Close() error
}

type RssDTO struct {
	Id     int    `json:"id"`
	Url    string `json:"url"`
	Rank   int    `json:"rank"`
	Title  string `json:"title"`
	Viewed bool   `json:"viewed"`
	Saved  bool   `json:"saved"`
}