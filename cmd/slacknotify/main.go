package main

import (
	"fmt"
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
	gatherMessages(c, l)

	// gets auth bearer token for twitch
	auth := gettoken(c, l)
	t := getStreamInfo(c, l, auth)
	fmt.Println(t)
}
