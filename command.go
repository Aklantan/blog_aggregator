package main

import (
	"fmt"
	"os"

	config "github.com/aklantan/blog_aggregator/internal"
)

type state struct {
	configuration *config.Config
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
	s.configuration.Current_user = cmd.arguments[0]
	err := config.WriteConfig(s.configuration) // Make sure you're passing the config correctly
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Println("User has been set")

	return nil
}
