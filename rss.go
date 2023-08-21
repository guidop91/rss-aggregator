package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title         string    `xml:"title"`
		Link          string    `xml:"link"`
		Description   string    `xml:"description"`
		Generator     string    `xml:"generator"`
		Language      string    `xml:"language"`
		LastBuildDate string    `xml:"lastBuildDate"`
		Items         []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	var emptyRSSFeed RSSFeed
	resp, fetchErr := httpClient.Get(url)
	if fetchErr != nil {
		// Return and log error upstream
		return emptyRSSFeed, fetchErr
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		// Return and log error upstream
		return emptyRSSFeed, readErr
	}

	rssFeed := RSSFeed{}

	unmarschalErr := xml.Unmarshal(body, &rssFeed)
	if unmarschalErr != nil {
		return emptyRSSFeed, unmarschalErr
	}

	return rssFeed, nil
}
