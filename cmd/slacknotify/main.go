package main

import (
	"fmt"
	"log"
	"os"
)

// example live stream info
var example []slackStreamInfo = []slackStreamInfo{{
	Name:   "b0aty",
	Title:  "this is rubbish",
	GameID: "459931",
	Link:   "",
}}

var exampleList slackStreamInfoList = slackStreamInfoList{
	list: example,
}

func main() {
	l := log.New(os.Stdout, "slack-notify-go", log.LstdFlags)

	// get config
	c, err := loadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		l.Fatal("Error loading config", err)
	}

	// sets up messages channel and goroutine to accept messages
	gatherMessages(c, l)

	// gets auth bearer token for twitch
	auth := gettoken(c, l)
	t := getStreamInfo(c, l, auth)
	fmt.Println(t)
	d := determineStatus(c, l, t)
	fmt.Println(d)

	uniqueIDs := returnUniqueIDs(l, exampleList)
	fmt.Println(uniqueIDs)
	games := findGameName(c, l, auth, uniqueIDs)
	fmt.Println(games)
}
