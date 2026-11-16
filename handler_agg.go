package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/DMCDanteh/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return err
	}
	data, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	for _, item := range data.Channel.Item {
		var publishedAt sql.NullTime
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			publishedAt.Time = parsedTime
			publishedAt.Valid = true
		}
		params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			FeedID:    feed.ID,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: publishedAt,
		}
		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23505" {
					continue
				}
			}
			fmt.Println(err)
			continue
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Usage: agg <time_between_reqs>")
	}
	timing, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}
	if timing <= 0 {
		return fmt.Errorf("duration must be greater than zero")
	}
	ticker := time.NewTicker(timing)
	defer ticker.Stop()
	fmt.Printf("Collecting feeds every %v\n", timing)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println(err)
		}
	}
}
