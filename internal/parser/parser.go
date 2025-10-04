package parser

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/jackc/pgx/v5/pgxpool"
	"telegram-news-bot/internal/db"
)

//Single-source parser
func ParseSource(pool *pgxpool.Pool, sourceID int, url string, category string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return fmt.Errorf("faild to parse feed: %w", err)
	}

	for _, item := range feed.Items {
		published := time.Now()
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		}

		err := db.AddNews(pool, sourceID, item.Title, item.Link, published, category)
		if err != nil {
			fmt.Printf("Error adding news: %v\n", err)
		}
	}

	return nil
}

//Cycle through all sources
func ParseAllSources(pool *pgxpool.Pool) {
	sources, err := db.GetAllSources(pool)
	if err != nil {
		fmt.Printf("Error getting sources: %v\n", err)
		return
	}

	for _, s := range sources {
		err := ParseSource(pool, s.ID, s.URL, s.Category)
		if err != nil {
			fmt.Printf("Error parsing source %s: %v\n", s.Name, err)
		}
	}
}
