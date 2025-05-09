package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	config "github.com/aklantan/blog_aggregator/internal"
	"github.com/aklantan/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	configuration *config.Config
	db            *database.Queries
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("state is nil: no configuration available")
	}
	handlerFunction, exists := c.commandList[cmd.name]
	if !exists {
		return fmt.Errorf("%s not found in command list", cmd.name)
	}

	return handlerFunction(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("no username provided")
		os.Exit(1)
		return fmt.Errorf("no username provided")
	}
	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Printf("% v", err)
		os.Exit(1)
	}
	s.configuration.Current_user = cmd.arguments[0]
	err = config.WriteConfig(s.configuration) // Make sure you're passing the config correctly
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Println("User has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("new user requires a name")
		os.Exit(1)
		return fmt.Errorf("no username provided")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.arguments[0]})
	if err != nil {
		fmt.Printf("User already exists %v\n", err)
		os.Exit(1)
	}
	s.configuration.Current_user = user.Name
	fmt.Printf(" %v %v %v %v\n", user.ID, user.Name, user.CreatedAt, user.UpdatedAt)
	return nil

}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("Cannot reset users table : %v", err)
	}
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Unable to list users : %v\n", err)
	}

	for _, user := range users {
		if user.Name == s.configuration.Current_user {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil

}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("must provide time between reqs")
		os.Exit(1)
	}
	frequency, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		fmt.Printf("cannot parse as a time")
		os.Exit(1)
	}

	fmt.Printf("Collecting feeds every %v\n", frequency)

	ticker := time.NewTicker(frequency)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.arguments) < 2 {
		fmt.Println("must provide name and Url")
		os.Exit(1)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.arguments[0], Url: cmd.arguments[1], UserID: user.ID})
	if err != nil {
		fmt.Printf("Cannot add feed to DB : %v", err)
		os.Exit(1)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		fmt.Printf("cannot add feed to follow list for user : %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(" %v %v %v %v %v %v\n", feed.ID, feed.Name, feed.CreatedAt, feed.UpdatedAt, feed.Url, feed.UserID)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("cannot retrieve feeds %v", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetFeedUser(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("cannot get feed creator : %v", err)
		}
		fmt.Printf("user : %v\n", user)
		fmt.Printf("name : %v\n url : %v\n", feed.Name, feed.Url)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		fmt.Println("url required")
		os.Exit(1)
		return fmt.Errorf("no username provided")
	}

	feed, err := s.db.GetFeedbyUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println("feed cannot be retrieved")
		os.Exit(1)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		fmt.Printf("cannot add feed to follow list for user : %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v %v\n", user.Name, feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Printf("cannot collect follows : %v\n", err)
		os.Exit(1)
	}

	for _, feed := range following {
		fmt.Printf("%v\n", feed.Name)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		fmt.Println("no username provided")
		os.Exit(1)
		return fmt.Errorf("no username provided")
	}
	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{Name: user.Name, Url: cmd.arguments[0]})
	if err != nil {
		fmt.Printf("cannot delete follow : %v\n", err)
	}
	return nil

}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit string = ""
	if len(cmd.arguments) == 0 {
		limit = "2"
	} else {
		limit = cmd.arguments[0]
	}
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		return fmt.Errorf("cannot convert to int32 : %w", err)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(intLimit)})
	if err != nil {
		return fmt.Errorf("cannot retrieve posts : %v", err)
	}

	for _, post := range posts {
		fmt.Printf("%v\n %v\n %v\n", post.Title, post.Url, post.Description)
	}

	return nil

}

//--------------- MIDDLEWARE -------------------------------

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.configuration.Current_user)
		if err != nil {
			fmt.Println("user cannot be retrieved")
			return err
		}
		return handler(s, cmd, user)
	}
}
