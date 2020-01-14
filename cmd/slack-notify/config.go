package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Twitch struct {
		API struct {
			Auth       string `json:"auth"`
			URLSTREAMS string `json:"url-streams"`
			URLUSERS   string `json:"url-users"`
		} `json:"api"`
		Streamers []string `json:"streamers"`
		Settings  struct {
			TIME string `json:"time"`
		} `json:"settings"`
	} `json:"twitch"`
	Slack struct {
		Webhook string `json:"webhook"`
	} `json:"Slack"`
}

func loadConfig(file string) (config, error) {
	var c config

	configFile, err := os.Open(file)

	defer configFile.Close()

	configJSON := json.NewDecoder(configFile)
	configJSON.Decode(&c)
	return c, err

}
