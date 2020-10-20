package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// base struct for a request
type twitchReq struct {
	url     string
	method  string
	headers []struct {
		key   string
		value string
	}
	params []string
}

// RESPONSE STRUCTS
// the structs below represent the json responses for each method
type twitchOAuthresp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// this is the base request for calls to the twitch api
func (t twitchReq) twitchRequest(l *log.Logger) []byte {
	var url = t.url

	// if params list isn't empty, add them to url
	if len(t.params) != 0 {
		url = fmt.Sprintf("%v%s", url, strings.Join(t.params, "&"))
	}

	req, err := http.NewRequest(t.method, url, nil)
	if err != nil {
		l.Println("Error sending twitchhttp request", err)
	}

	// add any headers present
	if len(t.headers) != 0 {
		for i := 0; i < len(t.headers); i++ {
			req.Header.Add(t.headers[i].key, t.headers[i].value)
		}
	}

	// send the request and return the response
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Println("Error reading response", err)
	}
	return respBody
}

// gets auth token from twitch
func gettoken(c config, l *log.Logger) string {
	r := &twitchReq{
		url:     c.Twitch.API.OAuthURL,
		method:  "POST",
		headers: nil,
		params: []string{
			c.Twitch.API.ClientID,
			c.Twitch.API.ClientSecret,
			"grant_type=client_credentials",
		},
	}
	response := r.twitchRequest(l)

	var oauthJSON twitchOAuthresp

	err := json.Unmarshal(response, &oauthJSON)
	if err != nil {
		l.Println("Error: could not parse token resp", err)
	}
	return oauthJSON.AccessToken
}
