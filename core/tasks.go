package ares

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

// TODO: provide better APIs for tasks
// tasks are one-time jobs where admins are expected it to run locally

// Removes all the existing images by going through each user list and
// re-uploads them to Imgur
// TODO: a lot of code duplicated
func (a *Ares) ClearImages() {
	api := slack.New(a.SlackAppToken)
	bot := slack.New(a.SlackBotToken)

	for {
		files, paging, err := api.GetFiles(slack.GetFilesParameters{
			Types: "images", Count: 1000})
		if err != nil {
			log.Fatal("Failed to fetch slack users info", err.Error())
		}
		if paging.Total <= 0 {
			log.Println("No more images to delete")
			return
		}
		for _, file := range files {
			resp := uploadToImgur(file.URLPrivateDownload, a.SlackAppToken, a.ImgurClientID)
			if resp.Status == false {
				log.Println("Failed to download/upload: ", file.ID)
				log.Println("Download url: ", file.URLPrivateDownload)
				continue
			}
			msg := fmt.Sprintf("Hi Deer, mesa Ares. I am cleaning up Slack and removing all the old images uploaded by you. Here is the imgur deletion link of below image: %s", resp.DeleteLink)
			params := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color:    "#D3D3D3",
				Text:     msg,
				ImageURL: resp.Link,
			}
			params.Attachments = []slack.Attachment{attachment}
			_, _, err := bot.PostMessage(file.User, "", params)
			if err != nil {
				log.Printf("Failed to send DM to %s: %s", file.User, err.Error())
				return
			}
			if err := api.DeleteFile(file.ID); err != nil {
				log.Printf("Deletion of the file %s failed: %s", file.URLPrivateDownload, err.Error())
				return
			}
			log.Println("Cleaned: ", file.ID)
		}
	}
}
