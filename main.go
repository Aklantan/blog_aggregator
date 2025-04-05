package main

import (
	"fmt"

	config "github.com/aklantan/blog_aggregator/internal"
)

func main() {
	configuration, _ := config.ReadConfig()
	configuration.SetUser()
	configuration, _ = config.ReadConfig()
	fmt.Println(configuration)

}
