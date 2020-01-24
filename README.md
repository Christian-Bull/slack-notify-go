# Slack Notify Go

This is a pretty simple app to post custom notifications to slack channels. As of now, this only pulls info from the twitch api. Notifying a slack channel when a user goes live. 


### Info

Put your config into `/etc/slack-notify/config.json` using the structure in the config example.


## Contents:

### `main.go`

Infinte loop that pings the slack api and determines if a user is live by using the `CreatedAt` field.


### `config.go`

Self explanatory, code to load in and parse the config file.


### `httpcalls.go`

Contains code for making a GET request on the twitch api and a POST request on the slack webhook. This should become a more robust package, but for now it's only a few lines of code and isn't difficult to re-use throughout the app.

### `twitchapi.go`

Parses the returned data from the twitch api and loads it into structs for use in other parts of the app.