package main

import (
	"fmt"
)

func main() {
	// loads config file into config struct
	c := loadConfig(`./assets/streamers.json`)

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
