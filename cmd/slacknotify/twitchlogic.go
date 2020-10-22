package main

import (
	"fmt"
	"log"
)

func printUsers(c config, l *log.Logger, a string) {
	t := getUserIDs(c, l, a)
	for i := 0; i < len(t.Data); i++ {
		fmt.Println(t.Data[i].Login, t.Data[i].ID)
	}
}
