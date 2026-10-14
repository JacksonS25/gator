package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/JacksonS25/gator/internal"
	"github.com/JacksonS25/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Type the command, not enought arguements provided.")
		os.Exit(1)
	}

	cfg, err := internal.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	state := internal.State{
		Config: &cfg,
		Db:     dbQueries,
	}

	commands := internal.Commands{
		Handlers: make(map[string]func(*internal.State, internal.Command) error),
	}

	commands.Register("login", internal.HandlerLogin)
	commands.Register("register", internal.HandlerRegister)
	commands.Register("reset", internal.HandlerReset)
	commands.Register("users", internal.HandlerUsers)
	commands.Register("agg", internal.HandlerAgg)
	commands.Register("addfeed", internal.MiddlewareLoggedIn(internal.HandlerAddFeed))
	commands.Register("feeds", internal.HandlerFeeds)
	commands.Register("follow", internal.MiddlewareLoggedIn(internal.HandlerFollow))
	commands.Register("following", internal.MiddlewareLoggedIn(internal.HandlerFollowing))
	commands.Register("unfollow", internal.MiddlewareLoggedIn(internal.HandlerUnfollow))
	commands.Register("browse", internal.MiddlewareLoggedIn(internal.HandlerBrowse))

	command := internal.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = commands.Run(&state, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
