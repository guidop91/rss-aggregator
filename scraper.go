package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/guidop91/rss-aggregator/internal/database"
)

const DateFormat = "Mon, 02 Jan 2006 15:04:05 -0700"

var existingPosts = make(map[string]bool)

func startScraping(
	db *database.Queries,
	concurrency int,
	requestInterval time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, requestInterval)

	ticker := time.NewTicker(requestInterval)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error fetching feeds", err)
			continue
		}

		urls, postsErr := db.GetPosts(context.Background())
		if postsErr != nil {
			log.Println("Error while fetching existing posts:", postsErr)
		}

		// create lookup table of existing post URLs
		for _, url := range urls {
			existingPosts[url] = true
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
		clearExistingPosts()
	}
}

func scrapeFeed(
	db *database.Queries,
	wg *sync.WaitGroup,
	feed database.Feed,
) {
	defer wg.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Could not mark feed as fetched", err)
		return
	}

	rssFeed, rssErr := urlToFeed(feed.Url)
	if rssErr != nil {
		log.Println("Error fetching feed", rssErr)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		// Find if post URL already exists in DB
		exists := existingPosts[item.Link]
		if exists {
			continue
		}

		// Parse pub date
		pubDate, dateErr := time.Parse(DateFormat, item.PubDate)
		if dateErr != nil {
			pubDate = time.Now().UTC()
		}

		// Parse description
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		_, postErr := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			Pubdate:     pubDate,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if postErr != nil {
			if strings.Contains(postErr.Error(), "duplicate key value") {
				continue
			}
			log.Println("Error creating post in db", postErr)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Items))
}

func clearExistingPosts() {
	existingPosts = make(map[string]bool)
}
