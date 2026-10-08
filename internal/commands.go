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
