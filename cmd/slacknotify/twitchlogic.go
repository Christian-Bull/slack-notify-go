package main

import (
	"fmt"
	"log"
)

// keeps all the info we need nice and tidy
type streamer struct {
	Name         string
	ID           string
	slackChannel string
}

// data needed for a message
type slackStreamInfo struct {
	Name string
	Game string
	Link string
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
