package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
		fmt.Printf("title : %v\n", item.Title)

	}
}
