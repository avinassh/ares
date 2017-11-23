package ares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
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

type VgyResponse struct {
	Error    bool   `json:"error"`
	URL      string `json:"url"`
	Image    string `json:"image"`
	Size     int    `json:"size"`
	Filename string `json:"filename"`
	Ext      string `json:"ext"`
	Delete   string `json:"delete"`
}

type UploadResponse struct {
	Link       string
	DeleteLink string
	Status     bool
}

func uploadToImgur(fileURL, slackAccessToken, imgurClientID string) *UploadResponse {
	result := &UploadResponse{}
	req, err := http.NewRequest("GET", fileURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slackAccessToken))
	client := &http.Client{
		Timeout: 3 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to download file from Slack: ", err.Error())
		log.Println("File url: ", fileURL)
		return result
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// not a valid response
		log.Println("Failed to download file from Slack:", resp.StatusCode)
		return result
	}

	// currently for some reason uploads to imgur are failing. following function first tries
	// to upload to Imgur and if that fails it tries to upload to Vgy
	imgBody, err := ioutil.ReadAll(resp.Body)
	imgurURL := "https://api.imgur.com/3/image"
	imgReq, err := http.NewRequest("POST", imgurURL, bytes.NewBuffer(imgBody))
	imgReq.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", imgurClientID))
	imgResp, err := client.Do(imgReq)
	if err != nil {
		log.Fatal("Failed to connect to Imgur: ", err.Error())
	}
	defer imgResp.Body.Close()
	if imgResp.StatusCode == 200 {
		var imgurResponse *ImgurResponse
		if err = json.NewDecoder(imgResp.Body).Decode(&imgurResponse); err != nil {
			log.Println("Failed to decode response from Imgur API: ", err.Error())
		}
		return formatImgurResponse(imgurResponse)

	}
	log.Println("Received a non-200 status while uploading to Imgur: ", imgResp.StatusCode)
	imgRespBody, _ := ioutil.ReadAll(imgResp.Body)
	log.Println("Imgur resp body for fail:", string(imgRespBody))
	// trying vgy
	vgyURL := "https://vgy.me/upload"
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", getFileNameFromSlackURL(fileURL))
	io.Copy(fw, bytes.NewReader(imgBody))
	w.Close()

	imgReq, err = http.NewRequest("POST", vgyURL, &b)
	imgReq.Header.Set("Content-Type", w.FormDataContentType())

	vgResp, err := client.Do(imgReq)
	defer vgResp.Body.Close()
	if vgResp.StatusCode != 200 {
		log.Println("Received a non-200 status while uploading to Vgy: ", vgResp.StatusCode)
		vgRespBody, _ := ioutil.ReadAll(vgResp.Body)
		log.Println("Vgy resp body for fail:", string(vgRespBody))
		return result
	}
	var vgResult *VgyResponse
	if err = json.NewDecoder(vgResp.Body).Decode(&vgResult); err != nil {
		log.Println("Failed to decode Vgy from Imgur API: ", err.Error())
		return result
	}
	return formatVgyResponse(vgResult)
}

// Slack URLs are of the format:
// https://files.slack.com/files-pri/T06V-F84/download/ggwp.png
// this function returns `ggwp.png`
func getFileNameFromSlackURL(url string) string {
	s := strings.Split(url, "/")
	return s[len(s)-1]
}

func formatImgurResponse(response *ImgurResponse) *UploadResponse {
	return &UploadResponse{
		Status:     true,
		Link:       response.Data.Link,
		DeleteLink: fmt.Sprintf("https://imgur.com/delete/%s", response.Data.Deletehash),
	}
}

func formatVgyResponse(response *VgyResponse) *UploadResponse {
	return &UploadResponse{
		Status:     true,
		Link:       response.Image,
		DeleteLink: response.Delete,
	}
}
