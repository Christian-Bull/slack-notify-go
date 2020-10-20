package main

import (
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "slack-notify-go", log.LstdFlags)

	// get config
	c, err := loadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		l.Fatal("Error loading config", err)
	}

	// sets up messages channel and goroutine to accept messages
	m := make(chan Message)
	go func() {
		err := sendMessages(c, l, m)
		if err != nil {
			l.Fatal(err)
		}
	}()

}
