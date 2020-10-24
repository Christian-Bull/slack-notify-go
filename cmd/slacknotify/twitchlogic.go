package main

import (
	"fmt"
	"log"
	"time"
)

var streamURL string = "https://www.twitch.tv"

// keeps all the info we need nice and tidy
type streamer struct {
	Name         string
	ID           string
	slackChannel string
}

// data needed for a message
type slackStreamInfo struct {
	Name   string
	Title  string
	GameID string // ugh idk why they return this as a string
	Link   string
}

type slackStreamInfoList struct {
	list []slackStreamInfo
}

// useful for testing, keeping it for now
func printUsers(c config, l *log.Logger, a string) {
	t := getUserIDs(c, l, a)
	for i := 0; i < len(t.Data); i++ {
		fmt.Println(t.Data[i].Login, t.Data[i].ID)
	}
}

// puts all the necessary info into a list of structs (id,name,slackchannel)
func (t twitchUserData) usersToList(c config, l *log.Logger) []streamer {
	var streamers []streamer

	for i := 0; i < len(t.Data); i++ {
		// find channel name from config
		var channel string
		for j := range c.Twitch.Streamers {
			if c.Twitch.Streamers[j].Name == t.Data[i].Login {
				channel = c.Twitch.Streamers[j].Channel
			}
		}
		// create a streamer struct and put into our array
		var s = streamer{
			Name:         t.Data[i].Login,
			ID:           t.Data[i].ID,
			slackChannel: channel,
		}
		streamers = append(streamers, s)
	}
	return streamers
}

func getStreamInfo(c config, l *log.Logger, auth string) livestreamers {
	t := getUserIDs(c, l, auth)
	streamers := t.usersToList(c, l)
	return getlivestreams(c, l, streamers, auth)
}

// this determines which streamers to send notifications for based on StartedAt
func determineStatus(c config, l *log.Logger, streams livestreamers) slackStreamInfoList {
	var liveStreams []slackStreamInfo

	for i := 0; i < len(streams.Data); i++ {
		now := time.Now().UTC()
		x, _ := time.ParseDuration(c.Twitch.Settings.Time)
		nowminus := now.Add(-x)

		// checks if started at time is after current time
		if streams.Data[i].StartedAt.After(nowminus) {
			data := slackStreamInfo{
				Name:   streams.Data[i].UserName,
				Title:  streams.Data[i].Title,
				GameID: streams.Data[i].GameID,
				Link:   getStreamURL(streams.Data[i].UserName),
			}
			liveStreams = append(liveStreams, data)
		}
	}
	var slackStreamers = slackStreamInfoList{
		list: liveStreams,
	}
	return slackStreamers
}

// returns streamer url
func getStreamURL(s string) string {
	return fmt.Sprintf("%s/%s", streamURL, s)
}

// replace game id with name
func (slack *slackStreamInfo) updateGame(n string) {
	slack.GameID = n
}

// returns unique games
func (s slackStreamInfoList) returnUniqueIDs(l *log.Logger) []string {
	var returnList []string
	for i := 0; i < len(s.list); i++ {
		for j := range returnList {
			if s.list[i].GameID == returnList[j] {
				continue
			} else {
				returnList = append(returnList, s.list[i].GameID)
			}
		}
	}
	return returnList
}

type game struct {
	ID   string
	Name string
}

// associates games and IDs
type games struct {
	games []game
}

// finds game names for a list of game ids
func findGameName(c config, l *log.Logger, auth string, gameIDs []string) games {
	data := getGameName(c, l, auth, gameIDs)

	var g []game

	for i := 0; i < len(data.Data); i++ {
		d := &game{
			ID:   data.Data[i].ID,
			Name: data.Data[i].Name,
		}
		g = append(g, *d)
	}
	return games{games: g}
}
