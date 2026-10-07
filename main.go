package main

import (
	"fmt"
	"os"

	"github.com/JacksonS25/gator/internal"
)

func main() {
	cfg, err := internal.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
	}

	state := internal.State{
		Config: &cfg,
	}

	commands := internal.Commands{
		Handlers: make(map[string]func(*internal.State, internal.Command) error),
	}

	commands.Register("login", internal.HandlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Type the command, not enought arguements provided.")
		os.Exit(1)
	}

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
