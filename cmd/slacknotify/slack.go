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
			l.Println("Error posting message: retry:", retries, err)
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
