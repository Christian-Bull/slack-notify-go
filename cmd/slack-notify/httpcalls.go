package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	// fmt.Println(resp.Status)

	return respBody
}

// takes (webhook, msg) and makes a post http req
func postMsg(webhook string, msg string) string {
	body := fmt.Sprintf(`{"text":"%s"}`, msg)

	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer([]byte(body)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.Status
}
