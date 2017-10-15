package ares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func uploadToImgur(fileURL, slackAccessToken, imgurClientID string) *ImgurResponse {
	var result *ImgurResponse
	req, err := http.NewRequest("GET", fileURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slackAccessToken))
	client := &http.Client{
		Timeout: 3 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to download file from Slack", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// not a valid response
		log.Println("Failed to download file from Slack:", resp.StatusCode)
		return result
	}
	imgBody, err := ioutil.ReadAll(resp.Body)
	imgURL := "https://api.imgur.com/3/image"

	imgReq, err := http.NewRequest("POST", imgURL, bytes.NewBuffer(imgBody))
	imgReq.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", imgurClientID))

	imgResp, err := client.Do(imgReq)

	if err != nil {
		log.Fatal("Failed to connect to Imgur", err.Error())
	}
	defer imgResp.Body.Close()
	if imgResp.StatusCode != 200 {
		// not a valid response
		log.Fatal("Received a non-200 status while uploading to Imgur", imgResp.StatusCode)
	}

	if err = json.NewDecoder(imgResp.Body).Decode(&result); err != nil {
		log.Fatal("Failed to decode response from Imgur API", err.Error())
	}

	return result
}
