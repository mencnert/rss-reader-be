package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Connector struct{}

func (c Connector) Fetch(url string) (RssFeed, error) {
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
