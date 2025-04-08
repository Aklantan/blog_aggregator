package main

import "fmt"

type state struct {
	configuration *Config
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no username provided")
	}
	s.configuration.Current_user = cmd.arguments[0]
	fmt.Println("User has been set")
	return nil
}
