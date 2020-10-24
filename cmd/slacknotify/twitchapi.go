package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var gameURL string = "https://api.twitch.tv/helix/games?"

// headers applied to an http request
type header struct {
	key   string
	value string
}

// parameters
type parameter struct {
	key   string
	value string
}

// base struct for a request
type twitchReq struct {
	url     string
	method  string
	headers []header
	params  []parameter
}

// this is the base request for calls to the twitch api
func (t twitchReq) twitchRequest(l *log.Logger) []byte {
	var url = t.url

	// if params list isn't empty, add them to url
	if len(t.params) != 0 {
		var params []string

		for i := 0; i < len(t.params); i++ {
			params = append(params, fmt.Sprintf("%s=%s", t.params[i].key, t.params[i].value))
		}
		url = fmt.Sprintf("%v%s", url, strings.Join(params, "&"))
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

// the structs below represent the json responses for each method
type twitchOAuthresp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// gets auth token from twitch
func gettoken(c config, l *log.Logger) string {
	var params = []parameter{
		parameter{
			key:   "client_id",
			value: c.Twitch.API.ClientID,
		},
		parameter{
			key:   "client_secret",
			value: c.Twitch.API.ClientSecret,
		},
		parameter{
			key:   "grant_type",
			value: "client_credentials",
		},
	}
	r := &twitchReq{
		url:     c.Twitch.API.OAuthURL,
		method:  "POST",
		headers: nil,
		params:  params,
	}
	response := r.twitchRequest(l)

	var oauthJSON twitchOAuthresp

	err := json.Unmarshal(response, &oauthJSON)
	if err != nil {
		l.Println("Error: could not parse token resp", err)
	}
	return oauthJSON.AccessToken
}

// list of user data
type twitchUserData struct {
	Data []struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		Type            string `json:"type"`
		BroadcasterType string `json:"broadcaster_type"`
		Description     string `json:"description"`
		ProfileImageURL string `json:"profile_image_url"`
		OfflineImageURL string `json:"offline_image_url"`
		ViewCount       int    `json:"view_count"`
	} `json:"data"`
}

func getUserIDs(c config, l *log.Logger, auth string) twitchUserData {
	// get the list of streamers and create the paramaters
	var params []parameter
	streamers := c.streamersToList()
	for i := 0; i < len(streamers); i++ {
		params = append(params, parameter{
			key:   "login",
			value: streamers[i],
		})
	}

	var requestHeaders = []header{
		header{
			key:   "Client-ID",
			value: c.Twitch.API.ClientID,
		},
		header{
			key:   "Authorization",
			value: fmt.Sprintf("Bearer %s", auth),
		},
	}

	// build the request
	r := &twitchReq{
		url:     c.Twitch.API.UserIDURL,
		method:  "GET",
		headers: requestHeaders,
		params:  params,
	}

	response := r.twitchRequest(l)

	var usersJSON twitchUserData

	err := json.Unmarshal(response, &usersJSON)
	if err != nil {
		l.Println("Error: could not parse user resp", err)
	}
	return usersJSON
}

// stream infomation struct
type livestreamers struct {
	Data []struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		UserName     string    `json:"user_name"`
		GameID       string    `json:"game_id"`
		Type         string    `json:"type"`
		Title        string    `json:"title"`
		ViewerCount  int       `json:"viewer_count"`
		StartedAt    time.Time `json:"started_at"`
		Language     string    `json:"language"`
		ThumbnailURL string    `json:"thumbnail_url"`
		TagIds       []string  `json:"tag_ids"`
	} `json:"data"`
	Pagination struct {
	} `json:"pagination"`
}

func getlivestreams(c config, l *log.Logger, streamers []streamer, auth string) livestreamers {
	// gets info from streamers that are live
	var params []parameter
	for i := 0; i < len(streamers); i++ {
		params = append(params, parameter{
			key:   "user_id",
			value: streamers[i].ID,
		})
	}

	// request header
	var requestHeaders = []header{
		header{
			key:   "Client-ID",
			value: c.Twitch.API.ClientID,
		},
		header{
			key:   "Authorization",
			value: fmt.Sprintf("Bearer %s", auth),
		},
	}

	// build the request
	r := &twitchReq{
		url:     c.Twitch.API.StreamInfoURL,
		method:  "GET",
		headers: requestHeaders,
		params:  params,
	}

	response := r.twitchRequest(l)

	var liveStreams livestreamers

	err := json.Unmarshal(response, &liveStreams)
	if err != nil {
		l.Println("Error: could not parse streamer info", err)
	}
	return liveStreams
}

// twitch api calls return gameid, we need to find the name
type gameResponse struct {
	Data []struct {
		BoxArtURL string `json:"box_art_url"`
		ID        string `json:"id"`
		Name      string `json:"name"`
	} `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

func getGameName(c config, l *log.Logger, auth string, gameIDs []string) gameResponse {
	var params []parameter
	for i := 0; i < len(gameIDs); i++ {
		params = append(params, parameter{
			key:   "id",
			value: gameIDs[i],
		})
	}

	var requestHeaders = []header{
		header{
			key:   "Client-ID",
			value: c.Twitch.API.ClientID,
		},
		header{
			key:   "Authorization",
			value: fmt.Sprintf("Bearer %s", auth),
		},
	}

	// build the request
	r := &twitchReq{
		url:     gameURL,
		method:  "GET",
		headers: requestHeaders,
		params:  params,
	}

	response := r.twitchRequest(l)

	var gamesJSON gameResponse

	err := json.Unmarshal(response, &gamesJSON)
	if err != nil {
		l.Println("Error: could not parse games response", err)
	}
	return gamesJSON
}
