package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Kristian-Roopnarine/rss/internal/database"
)

func Scraper(db *database.Queries, concurrency int, waitTime time.Duration) {
	ticker := time.NewTicker(waitTime)
	log.Printf("Collecting feeds every %s on %v goroutines...\n", waitTime, concurrency)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch")
			continue
		}

		log.Printf("Found %v feeds to fetch", len(feeds))
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		NewTime: time.Now().UTC(),
		FeedID:  feed.ID,
	})
	if err != nil {
		log.Printf("Couldn't mark feed %s as fetched: %v", feed.Name, err)
		return
	}

	feedData, err := FetchRSSFromURL(feed.Url)
	if err != nil {
		log.Printf("Couldn't fetch feed for %s: %v", feed.Url, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		log.Println("Found post", item.Title)
	}
}
