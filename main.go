package main

import (
	"fmt"
	"os"

	config "github.com/aklantan/blog_aggregator/internal"
)

func main() {
	configuration, err := config.ReadConfig()
	if err != nil {
		fmt.Println("no config file in config location")
	}
	appState := state{
		configuration: &configuration,
	}

	commands := commands{
		commandList: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("no command provided")
		os.Exit(1)
	}

	newCommand := command{
		name:      args[1],
		arguments: args[2:],
	}

	commands.run(&appState, newCommand)

}
