package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Twitch struct {
		API struct {
			ClientID      string `json:"client_id"`
			ClientSecret  string `json:"client_secret"`
			OAuthURL      string `json:"OAuthURL"`
			UserIDURL     string `json:"UserIDURL"`
			StreamInfoURL string `json:"StreamInfoURL"`
		} `json:"api"`
		Streamers []struct {
			Name    string `json:"name"`
			Channel string `json:"channel"`
		} `json:"streamers"`
		Settings struct {
			Time string `json:"time"`
		} `json:"settings"`
	} `json:"twitch"`
	Slack struct {
		Webhook     string `json:"webhook"`
		Auth        string `json:"auth"`
		Log         string `json:"log"`
		Postchannel string `json:"postchannel"`
		Logchannel  string `json:"logchannel"`
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

func (c config) streamersToList() []string {
	var streamers []string
	for i := 0; i < len(c.Twitch.Streamers); i++ {
		streamers = append(streamers, c.Twitch.Streamers[i].Name)
	}
	return streamers
}
