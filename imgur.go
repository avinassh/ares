package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ImgurResponse struct {
	Data struct {
		Deletehash string `json:"deletehash"`
		Link       string `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

func uploadToImgur(fileURL, slackAccessToken string) {
	clientID := os.Getenv("IMGUR_CLIENT_ID")

	req, err := http.NewRequest("GET", fileURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slackAccessToken))
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to connect to Slack", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// not a valid response
		log.Println(resp.StatusCode)
		return
	}

	iurl := "https://api.imgur.com/3/image"

	ireq, err := http.NewRequest("POST", iurl, resp.Body)
	ireq.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", clientID))

	iresp, err := client.Do(ireq)

	if err != nil {
		log.Fatal("Failed to connect to Imgur", err.Error())
	}
	defer iresp.Body.Close()
	if iresp.StatusCode != 200 {
		// not a valid response
		log.Println(resp.StatusCode)
		return
	}

	var result ImgurResponse
	if err = json.NewDecoder(iresp.Body).Decode(&result); err != nil {
		log.Fatal("Failed to decode response from Chat API", err.Error())
	}

	log.Println(result.Data.Link, result.Data.Deletehash)

}
