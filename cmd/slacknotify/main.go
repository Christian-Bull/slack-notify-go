package main

import (
	"log"
	"os"
	"time"
)

func main() {
	l := log.New(os.Stdout, "slack-notify-go", log.LstdFlags)

	// get config
	c, err := loadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		l.Fatal("Error loading config", err)
	}

	// sets up our send message channel
	m := gatherMessages(c, l)

	// gets auth bearer token for twitch
	auth := gettoken(c, l)

	for {
		// runs bot
		runTwitchBot(c, l, auth, m)

		// sleep for time in config
		sleepTime, err := time.ParseDuration(c.Twitch.Settings.Time)
		l.Printf("Waiting %s before next http call", sleepTime)
		if err != nil {
			l.Fatal("Couldn't parse sleep time", err)
		} else {
			time.Sleep(sleepTime)
		}
	}

}
