package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Twitch struct {
		API struct {
			Auth       string `json:"auth"`
			URLStreams string `json:"url-streams"`
			URLUsers   string `json:"url-users"`
		} `json:"api"`
		Streamers []string `json:"streamers"`
		Settings  struct {
			Time string `json:"time"`
		} `json:"settings"`
	} `json:"twitch"`
	Slack struct {
		Webhook string `json:"webhook"`
	} `json:"slack"`
}

func loadConfig(file string) (config, error) {
	var c config

	configFile, err := os.Open(file)

	defer configFile.Close()

	configJSON := json.NewDecoder(configFile)
	configJSON.Decode(&c)
	return c, err
}
