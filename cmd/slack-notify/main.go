package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// initiate logging to stdout
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "Logger: ", 2)
	)

	logger.SetOutput(os.Stdout)
	logger.Print("Logging started")

	// loads config file into config struct
	c, err := loadConfig(`/etc/slack-notify/config.json`)
	if err != nil {
		logger.Printf("error loading config: %v", err)
		logger.Printf("Exit")
		os.Exit(1)
	}

	logger.Printf(c.Slack.Webhook)
	slackstatus := postMsg(c.Slack.Webhook, "Connected")
	logger.Printf("Slack connection: %s", slackstatus)

	// infinite loop
	for {
		// gets twitch data
		logger.Print("Getting IDs")
		ch, _, s, _ := twitch(c)
		for i := 0; i < ch.Total; i++ {
			logger.Printf("Getting stream data for %s", ch.Users[i].Name)
		}

		// finds live streams and posts msgs
		live := twitchLive(c, s)
		if live != nil {
			slackPost(c, live)

		}

		// sleep for time in config
		t, err := time.ParseDuration(c.Twitch.Settings.Time)
		logger.Printf("Waiting %s before next http call", t)
		if err != nil {
			logger.Panic(err)
		} else {
			time.Sleep(t)
		}
	}
}

// channel struct, IDs, livedata, offline streams
func twitch(c config) (channelName, []string, []streamData, []offline) {
	channelData, IDs := c.getIDs()

	streamData, offline := c.getStreamData(IDs)

	return channelData, IDs, streamData, offline // for other uses
}

func twitchLive(c config, s []streamData) []streamData {
	var channels []streamData

	log.Printf("Live streams:")
	// only for online channels
	for i := 0; i < len(s); i++ {
		now := time.Now().UTC()
		x, _ := time.ParseDuration(c.Twitch.Settings.Time)
		nowminus := now.Add(-x)

		if s[i].Stream.CreatedAt.After(nowminus) {
			channels = append(channels, s[i])
		} else {
			fmt.Println(s[i].Stream.Channel.Name, s[i].Stream.CreatedAt)
		}
	}
	return channels // returns live users
}

func slackPost(c config, s []streamData) {
	// url := c.Slack.Webhook

	for i := 0; i < len(s); i++ {
		name := s[i].Stream.Channel.Name
		game := s[i].Stream.Game
		status := s[i].Stream.Channel.Status
		link := s[i].Stream.Channel.URL

		msg := fmt.Sprintf("%s is now live! Game: %s\n`%s`\n%s", name, game, status, link)
		_ = postMsg(c.Slack.Webhook, msg)
		log.Printf(msg)
	}
}
