package main

import (
	"fmt"
	"time"
)

func main() {
	// loads config file into config struct
	c := loadConfig(`../../assets/config.json`)

	_, IDs := c.getIDs()
	fmt.Println(IDs)

	streamData, _ := c.getStreamData(IDs)

	for i := 0; i < len(streamData); i++ {
		streamData[i].print()
		time := streamData[i].Stream.CreatedAt
		// diff := now.Sub(time).Seconds
		fmt.Println(time)
	}

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

// func test() {
// 	streamers := []string(strings.Split(readFile("assets/streamers.txt"), "\n"))

// 	_, IDs := getIDs(streamers)

// 	//fmt.Println(channelnames)

// 	streamData, _ := getStreamData(IDs)

// 	// x, _ := time.ParseDuration("5m")
// 	now := time.Now().UTC()
// 	fmt.Println(now)

// 	// for i := 0; i < len(streamData); i++ {
// 	// 	fmt.Println(streamData[i].Stream.CreatedAt)
// 	// }
// 	// for i := 0; i < len(offlineData); i++ {
// 	// 	fmt.Println(offlineData[i])

// 	for i := 0; i < len(streamData); i++ {
// 		streamData[i].print()
// 		time := streamData[i].Stream.CreatedAt
// 		// diff := now.Sub(time).Seconds
// 		fmt.Println(time)
// 	}
// }
