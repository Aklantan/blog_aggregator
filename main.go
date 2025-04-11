package main

import (
	"database/sql"
	"fmt"
	"os"

	config "github.com/aklantan/blog_aggregator/internal"
	"github.com/aklantan/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	configuration, err := config.ReadConfig()
	if err != nil {
		fmt.Println("no config file in config location")
	}
	db, err := sql.Open("postgres", configuration.Db_url)
	if err != nil {
		fmt.Println("cannot open database")
		os.Exit(1)
	}
	dbQueries := database.New(db)

	appState := state{
		configuration: &configuration,
		db:            dbQueries,
	}

	commands := commands{
		commandList: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)

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
	config.WriteConfig(&configuration)

}
