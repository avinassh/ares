package ares

import (
	"fmt"

	"github.com/nlopes/slack"
	"log"
)

type Ares struct {
	SlackAppToken string
	SlackBotToken string
	ImgurClientID string
}

func (a *Ares) deleteFile(fileId string) {
	api := slack.New(a.SlackAppToken)
	if err := api.DeleteFile(fileId); err != nil {
		fmt.Println("Deletion of the file failed:", err.Error())
	}
}

func (a *Ares) notifyUser(user, deleteHash string) {
	api := slack.New(a.SlackBotToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color:   "danger",
		Pretext: "Hey there, image uploads to this slack is disabled. Please use external image upload services next time. This incident has been reported :squirrel:",
		Text:    "I have uploaded the image to Imgur for now.",
		Fields: []slack.AttachmentField{
			{
				Title: "If you wish to delete the image uploaded, click on the following link",
				Value: fmt.Sprintf("https://imgur.com/delete/%s", deleteHash),
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(user, "", params)
	if err != nil {
		log.Println("Failed to send DM to user", err.Error())
	}
}

func (a *Ares) sendImgToSlack(channel, user, url string) {
	api := slack.New(a.SlackBotToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color:    "good",
		Text:     fmt.Sprintf("image originally uploaded by <@%s>", user),
		ImageURL: url,
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(channel, "", params)
	if err != nil {
		log.Println("Failed to send DM to user", err.Error())
	}
}

func (a *Ares) handleFile(file *slack.File, channel string) {
	resp := uploadToImgur(file.URLPrivateDownload, a.SlackAppToken, a.ImgurClientID)

	if resp.Status != 200 {
		log.Println("Failed to download/upload")
		return
	}

	a.notifyUser(file.User, resp.Data.Deletehash)
	a.sendImgToSlack(channel, file.User, resp.Data.Link)
	a.deleteFile(file.ID)
}

func isImageFile(fileType string) bool {
	for _, fType := range []string{"jpg", "jpeg", "png", "gif"} {
		if fType == fileType {
			return true
		}
	}
	return false
}

func (a *Ares) Run() {
	api := slack.New(a.SlackBotToken)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			if ev.SubType == "file_share" {
				if isImageFile(ev.File.Filetype) {
					a.handleFile(ev.File, ev.Channel)
				}
			}

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
