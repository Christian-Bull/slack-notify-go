package main

import (
	"log"

	"github.com/slack-go/slack"
)

// Message is the struct used to format a slack message
type Message struct {
	message string
	channel string
	status  string
}

// gets messages from our m channel and submits a post request
func sendMessages(c config, l *log.Logger, m chan Message) error {
	return postMessage(c, l, <-m)
}

// Post a message to slack
func postMessage(c config, l *log.Logger, m Message) error {
	var (
		retries int = 3
		err     error
	)

	api := slack.New(c.Slack.Auth)

	// retry slack post until it hits the retry limit or is successful
	for retries > 0 {
		msgID, _, _, err := api.SendMessage(
			m.channel,
			slack.MsgOptionText(m.message, false),
		)
		if err != nil {
			l.Println("Error posting message: retrying", retries)
			retries--
		} else {
			l.Println("Sent message to: ", m.channel, msgID)
			break
		}

	}
	return err
}

func createMessage(message string, channel string) Message {
	return Message{
		message: message,
		channel: channel,
		status:  "",
	}
}

func gatherMessages(c config, l *log.Logger) {
	// test our connection using the default post channel
	err := postMessage(c, l, createMessage("Connected", c.Slack.Logchannel))
	if err != nil {
		l.Fatal("Couldn't send initial message", err)
	}

	m := make(chan Message)
	go func() {
		err := sendMessages(c, l, m)
		if err != nil {
			l.Fatal(err)
		}
	}()
}
