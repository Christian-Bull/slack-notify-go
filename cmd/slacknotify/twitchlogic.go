package main

import (
	"fmt"
	"log"
	"strings"
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
	Name        string
	Title       string
	GameID      string // ugh idk why they return this as a string
	Link        string
	PostChannel string
}

type slackStreamInfoList struct {
	list []slackStreamInfo
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

// returns streamer url
func getStreamURL(s string) string {
	return fmt.Sprintf("%s/%s", streamURL, s)
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

// replace game id with name
func (slack *slackStreamInfo) updateGame(n string) {
	slack.GameID = n
}

// returns unique games
func returnUniqueIDs(l *log.Logger, s slackStreamInfoList) []string {
	// put our gameIDs in a list
	var gameIDs []string
	for j := 0; j < len(s.list); j++ {
		gameIDs = append(gameIDs, s.list[j].GameID)
	}

	// check if they're unique, append if not
	keys := make(map[string]bool)
	list := []string{}
	for _, i := range gameIDs {
		if _, value := keys[i]; !value {
			keys[i] = true
			list = append(list, i)
		}
	}
	return list
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

func (s *slackStreamInfoList) updateGameID(l *log.Logger, g games) {
	for i := 0; i < len(s.list); i++ {
		for j := 0; j < len(g.games); j++ {
			if s.list[i].GameID == g.games[j].ID {
				s.list[i].updateGame(g.games[j].Name)
			}
		}
	}
}

func (slack *slackStreamInfo) updatePostchannel(c string) {
	slack.PostChannel = c
}

// updates our livestream struct with post channel
func (s *slackStreamInfoList) updatePostChannels(c config, l *log.Logger) {
	for i := 0; i < len(s.list); i++ {
		for _, j := range c.Twitch.Streamers {
			if strings.ToLower(s.list[i].Name) == strings.ToLower(j.Name) {
				s.list[i].updatePostchannel(j.Channel)
			}
		}
	}
}

func (slack slackStreamInfo) formatMessage() string {
	name := slack.Name
	game := slack.GameID
	title := slack.Title
	link := slack.Link
	return fmt.Sprintf("%s is now live! Game: %s\n`%s`\n%s", name, game, title, link)
}

func (s *slackStreamInfoList) sendNotifications(c config, l *log.Logger, m chan Message) {
	for i := 0; i < len(s.list); i++ {
		msg := createMessage(s.list[i].formatMessage(), s.list[i].PostChannel)
		m <- msg
	}
}

func runTwitchBot(c config, l *log.Logger, auth string, m chan Message) {
	t := getStreamInfo(c, l, auth)
	for i := 0; i < len(t.Data); i++ {
		l.Printf("%s is live, started at %s", t.Data[i].UserName, t.Data[i].StartedAt.String())
	}

	d := determineStatus(c, l, t)

	uniqueIDs := returnUniqueIDs(l, d)
	games := findGameName(c, l, auth, uniqueIDs)
	d.updateGameID(l, games)
	d.updatePostChannels(c, l)
	if d.list == nil {
		l.Println("No new streams to post info for")
	} else {
		l.Println("Found the following streams to post a notification for", d.list)
	}
	d.sendNotifications(c, l, m)
}
