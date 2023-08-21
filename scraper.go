package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/guidop91/rss-aggregator/internal/database"
)

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

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
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
		log.Printf("On Feed \"%s\" , post with title \"%s\" found.\n", feed.Name, item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Items))
}
