package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type channelName struct {
	Total int `json:"_total"`
	Users []struct {
		DisplayName string    `json:"display_name"`
		ID          string    `json:"_id"`
		Name        string    `json:"name"`
		Type        string    `json:"type"`
		Bio         string    `json:"bio"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Logo        string    `json:"logo"`
	} `json:"users"`
}

type streamData struct {
	Stream struct {
		ID          int64     `json:"_id"`
		Game        string    `json:"game"`
		Viewers     int       `json:"viewers"`
		VideoHeight int       `json:"video_height"`
		AverageFps  int       `json:"average_fps"`
		Delay       int       `json:"delay"`
		CreatedAt   time.Time `json:"created_at"`
		IsPlaylist  bool      `json:"is_playlist"`
		Preview     struct {
			Small    string `json:"small"`
			Medium   string `json:"medium"`
			Large    string `json:"large"`
			Template string `json:"template"`
		} `json:"preview"`
		Channel struct {
			Mature                       bool        `json:"mature"`
			Status                       string      `json:"status"`
			BroadcasterLanguage          string      `json:"broadcaster_language"`
			DisplayName                  string      `json:"display_name"`
			Game                         string      `json:"game"`
			Language                     string      `json:"language"`
			ID                           int         `json:"_id"`
			Name                         string      `json:"name"`
			CreatedAt                    time.Time   `json:"created_at"`
			UpdatedAt                    time.Time   `json:"updated_at"`
			Partner                      bool        `json:"partner"`
			Logo                         string      `json:"logo"`
			VideoBanner                  string      `json:"video_banner"`
			ProfileBanner                string      `json:"profile_banner"`
			ProfileBannerBackgroundColor interface{} `json:"profile_banner_background_color"`
			URL                          string      `json:"url"`
			Views                        int         `json:"views"`
			Followers                    int         `json:"followers"`
		} `json:"channel"`
	} `json:"stream"`
}

type offline struct {
	Stream string `json:"stream"`
}

func readFile(f string) string {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	return string(b)
}

// TODO
// need to change this to (interface?)
// make it able to return a variable struct type
//
// func ToStruct(b byte) interface{} {
// 	return interface
// }

func channelToStruct(j string) channelName {
	var channelJSON channelName

	err := json.Unmarshal([]byte(j), &channelJSON)
	if err != nil {
		fmt.Println(err)
	}

	return channelJSON
}

func streamToStruct(j string) streamData {
	var streamJSON streamData

	err := json.Unmarshal([]byte(j), &streamJSON)
	if err != nil {
		fmt.Println(err)
	}
	return streamJSON
}

func offlineToStruct(j string) (offline, error) {
	var offlineJSON offline

	err := json.Unmarshal([]byte(j), &offlineJSON)

	return offlineJSON, err
}

func (c config) getIDs() (channelName, []string) {
	// puts streamers into a list, ugh this is a hotifx
	var u []string
	for i := 0; i < len(c.Twitch.Streamers); i++ {
		u = append(u, c.Twitch.Streamers[i].Name)
	}

	auth := c.Twitch.API.Auth
	url := c.Twitch.API.URLUsers

	channelStructs := channelToStruct(string(httpGet(url+strings.Join(u, ","), auth)))
	if channelStructs.Total == 0 {
		fmt.Println("No IDs found")
		os.Exit(1)
	}
	var UserIds []string
	for i := 0; i < channelStructs.Total; i++ {
		UserIds = append(UserIds, channelStructs.Users[i].ID)
	}
	return channelStructs, UserIds
}

func (c config) getStreamData(u []string) ([]streamData, []offline) {
	auth := c.Twitch.API.Auth
	url := c.Twitch.API.URLStreams

	var streamData []streamData
	var offlineData []offline

	// I feel like this is bad code
	// Checks if a channel is offline
	for i := 0; i < len(u); i++ {
		resp := string(httpGet(url+u[i], auth))
		off, err := offlineToStruct(resp)
		if off.Stream != "" || err != nil {
			item := streamToStruct(resp)
			streamData = append(streamData, item)
			fmt.Println("User: %s is live", item.Stream.Channel.Name)
		} else {
			offlineData = append(offlineData, off)
			fmt.Println("user is offline")
		}
	}

	return streamData, offlineData
}

func (s streamData) print() {
	fmt.Printf("%+v", s)
}
