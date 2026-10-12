package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JacksonS25/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Config *Config
	Db     *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	err := c.Handlers[cmd.Name](s, cmd)
	if err != nil {
		return fmt.Errorf("error running command %s: %w", cmd.Name, err)
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: login <username>")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("error logging in user: %v\n", err)
		os.Exit(1)
	}

	SetUser(*s.Config, cmd.Args[0])

	fmt.Println("User has been set: ", cmd.Args[0])

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: register <username>")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	user, err := s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Printf("error creating user: %v\n", err)
		os.Exit(1)
	}

	SetUser(*s.Config, cmd.Args[0])
	fmt.Println("User has been registered: ", cmd.Args[0])
	fmt.Println(user)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("error resetting users: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All users have been reset.")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error getting users: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Users:")
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Printf("error fetching feed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(feed.Channel.Title)
	fmt.Println(feed.Channel.Link)
	fmt.Println(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		fmt.Println("Title:", item.Title)
		fmt.Println("Link:", item.Link)
		fmt.Println("Description:", item.Description)
		fmt.Println("PubDate:", item.PubDate)
		fmt.Println()
	}
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		fmt.Printf("error getting current user: %v\n", err)
		os.Exit(1)
	}

	feedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUser.ID,
	}

	feed, err := s.Db.AddFeed(context.Background(), feedParams)
	if err != nil {
		fmt.Printf("error adding feed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Feed has been added: ", feed.Name)
	cmd.Args = cmd.Args[1:]
	err = HandlerFollow(s, cmd)
	if err != nil {
		fmt.Printf("error following feed: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.Feeds(context.Background())
	if err != nil {
		fmt.Printf("error getting feeds: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Feeds:")
	for _, feed := range feeds {
		user, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("error getting user: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("* %s (%s): %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: follow <feed_url>")
	}

	currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		fmt.Printf("error getting current user: %v\n", err)
		os.Exit(1)
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("error getting feed: %v\n", err)
		os.Exit(1)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		fmt.Printf("error following feed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("User %s is now following feed %s\n", currentUser.Name, feed.Name)
	return nil
}

func HandlerFollowing(s *State, cmd Command) error {
	currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		fmt.Printf("error getting current user: %v\n", err)
		os.Exit(1)
	}

	following, err := s.Db.GetFeedFollowForUser(context.Background(), currentUser.ID)
	if err != nil {
		fmt.Printf("error getting following feeds: %v\n", err)
		os.Exit(1)
	}

	feeds := []string{}
	for _, feedFollow := range following {
		feed, err := s.Db.GetFeedById(context.Background(), feedFollow.FeedID)
		if err != nil {
			fmt.Printf("error getting feed: %v\n", err)
			os.Exit(1)
		}
		feeds = append(feeds, feed.Name)
	}
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed)
	}
	return nil
}
