package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	l := log.New(os.Stdout, "slack-notify-go", log.LstdFlags)

	// get config
	c, err := loadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		l.Fatal("Error loading config", err)
	}

	// gets auth bearer token for twitch
	auth := gettoken(c, l)

	// test our connection using the default post channel
	err = postMessage(c, l, createMessage("Connected", c.Slack.Logchannel))
	if err != nil {
		l.Fatal("Couldn't send initial message", err)
	}

	go func() {
		for {
			// runs bot
			runTwitchBot(c, l, auth)

			// sleep for time in config
			sleepTime, err := time.ParseDuration(c.Twitch.Settings.Time)
			l.Printf("Waiting %s before next http call", sleepTime)
			if err != nil {
				l.Fatal("Couldn't parse sleep time", err)
			} else {
				time.Sleep(sleepTime)
			}
		}
	}()

	// basic prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("LISTENING_PORT")
	if port == "" {
		port = "8080"
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		l.Fatal("Failed to start prometheus metrics server")
	}
}
