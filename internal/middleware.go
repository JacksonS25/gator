package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/JacksonS25/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			fmt.Printf("error getting current user: %v\n", err)
			os.Exit(1)
		}
		return handler(s, cmd, currentUser)
	}
}
