package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
)

// GET request, takes url with parameters set
func httpGet(url string, auth string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Client-ID", auth)
	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(respBody)

	return respBody
}

func postMsg(auth string, channel string, msg string) (string, error) {
	api := slack.New(auth)

	msgID, _, _, err := api.SendMessage(
		channel,
		slack.MsgOptionText(msg, false),
	)

	return msgID, err
}
