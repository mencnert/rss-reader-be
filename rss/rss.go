package rss

import (
	"encoding/xml"
)

type RssFeed struct {
	XmlName xml.Name   `xml:"feed"`
	Entries []RssEntry `xml:"entry"`
}

type RssEntry struct {
	XmlName xml.Name `xml:"entry" json:"-"`
	Url     string   `xml:"id" json:"url"`
	Rank    int      `xml:"rank" json:"rank"`
	Title   string   `xml:"title" json:"title"`
}
