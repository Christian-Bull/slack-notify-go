package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Twitch struct {
		API struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			OAuthURL     string `json:"OAuthURL"`
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
	} `json:"slack"`
	Darksky struct {
		Secret string `json:"secret"`
		URL    string `json:"url"`
	} `json:"darksky"`
}

func loadConfig(file string) (config, error) {
	var c config

	configFile, err := os.Open(file)

	defer configFile.Close()

	configJSON := json.NewDecoder(configFile)
	configJSON.Decode(&c)
	return c, err
}
