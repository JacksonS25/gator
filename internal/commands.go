package internal

import (
	"fmt"
)

type State struct {
	Config *Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*State, command) error
}

func (c *commands) run(s *State, cmd command) error {
	err := c.handlers[cmd.name](s, cmd)
	if err != nil {
		return fmt.Errorf("error running command %s: %w", cmd.name, err)
	}
	return nil
}

func (c *commands) register(name string, f func(*State, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *State, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: login <username>")
	}

	SetUser(*s.Config, cmd.args[0])

	fmt.Println("User has been set: ", cmd.args[0])

	return nil
}
