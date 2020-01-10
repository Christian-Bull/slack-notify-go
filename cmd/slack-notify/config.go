package main

import (
	"encoding/json"
	"fmt"
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
	} `json:"twitch"`
}

func loadConfig(file string) config {
	var c config

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	configJSON := json.NewDecoder(configFile)
	configJSON.Decode(&c)
	return c

}
