package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "slack-notify-go", log.LstdFlags)

	// get config
	c, err := loadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		logger.Fatal("Error loading config", err)
	}

	// handle events
	err := 

	// create tasks
	err := runTasks()
	if err != nil {
		logger.Fatal("Error running tasks", err)
	}

	// run tasks

}
