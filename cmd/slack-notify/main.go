package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// initiate logging to stdout
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "Logger: ", 2)
	)
	logger.SetOutput(os.Stdout)
	logger.Print("Logging started")

	// loads config file into config struct
	c, err := loadConfig(`../../assets/streamers.json`)
	if err != nil {
		logger.Printf("error loading config: %v", err)
	}

	twitch(c)

}

func twitch(c config) {
	_, IDs := c.getIDs()

	streamData, _ := c.getStreamData(IDs)

	for i := 0; i < len(streamData); i++ {
		now := time.Now().UTC()
		x, _ := time.ParseDuration(c.Twitch.Settings.TIME)
		nowminus := now.Add(-x)

		if nowminus.After(streamData[i].Stream.CreatedAt) {
			fmt.Println(streamData[i].Stream.CreatedAt)
		} else {
			fmt.Println("Not After")
		}
	}
}
