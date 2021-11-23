package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type SofConnector interface {
	Get(string) (RssFeed, error)
}

type SOF struct{}

func (sof SOF) Get(url string) (RssFeed, error) {
	res, err := http.Get(url)
	if err != nil {
		return RssFeed{}, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return RssFeed{}, err
	}

	var feed RssFeed
	err = xml.Unmarshal(body, &feed)

	if err != nil {
		return RssFeed{}, err
	}

	return feed, nil
}
