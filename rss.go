package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/aklantan/blog_aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Printf("request is not possible : %v\n", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("response not possible : %v\n", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("response code is not 200 : %v\n", resp.StatusCode)
		return nil, err
	}
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("unable to read response : %v\n", err)
		return nil, err
	}
	feed := RSSFeed{}
	err = xml.Unmarshal(bodyByte, &feed)
	if err != nil {
		fmt.Printf("unable to unmarshal xml : %v\n", err)
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil

}

func scrapeFeeds(s *state) {
	nxtFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("cannot get next feed\n", err)
		return
	}
	err = s.db.MarkFeedFetched(context.Background(), nxtFeed.ID)
	if err != nil {
		fmt.Println("cannot mark feed fetched")
		return
	}
	fetchedFeed, err := fetchFeed(context.Background(), nxtFeed.Url)
	if err != nil {
		fmt.Println("cannot fetch feeds")
		return
	}
	for _, item := range fetchedFeed.Channel.Item {
		title := sql.NullString{String: item.Title, Valid: true}       // convert title to NullStr
		descr := sql.NullString{String: item.Description, Valid: true} // convert description to NullStr

		dateFormats := []string{
			"Mon, 02 Jan 2006 15:04:05 MST", // RSS 2.0 format
			"2006-01-02T15:04:05Z",          // ISO 8601
			"2006-01-02",                    // Simple YYYY-MM-DD
		}

		var parsedDate time.Time

		for _, format := range dateFormats {
			parsedDate, err = time.Parse(format, item.PubDate)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Println("unable to parse date")
		}

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: title, Description: descr, Url: item.Link, PublishedAt: parsedDate, FeedID: nxtFeed.ID})
		if err != nil {
			if err != nil {
				// Check if it's a PostgreSQL error
				if pqErr, ok := err.(*pq.Error); ok {
					// Inspect the code for duplicate key
					if pqErr.Code == "23505" { // Unique constraint violation code
						fmt.Println("Duplicate post URL, ignored.")
					}
				}
				// Log or handle other unexpected errors
				fmt.Printf("unexpected error: %v", err)
			}
		}
	}

}
