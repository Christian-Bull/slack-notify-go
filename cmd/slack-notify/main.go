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
	c, err := loadConfig(`../../assets/config.json`)
	if err != nil {
		logger.Printf("error loading config: %v", err)
	}

	// infinite loop
	for {
		// gets twitch data
		logger.Print("Getting IDs")
		_, IDs, s, o := twitch(c)
		logger.Printf("Getting stream data for %s", IDs)
		logger.Print(o)

		// finds live streams and posts msgs
		live := twitchLive(c, s)
		if live != nil {
			slackPost(c, live)
		}

		// sleep for time in config
		t, err := time.ParseDuration(c.Twitch.Settings.TIME)
		logger.Print(t)
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

	// only for online channels
	for i := 0; i < len(s); i++ {
		now := time.Now().UTC()
		x, _ := time.ParseDuration(c.Twitch.Settings.TIME)
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

		msg := fmt.Sprintf("%s is now live! Game: %s\n%s\n%s", name, game, status, link)
		//postMsg(url, msg)
		fmt.Println(msg)
	}
}
