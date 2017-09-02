package ares

import (
	"fmt"

	"github.com/nlopes/slack"
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

func (a *Ares) handleFile(file *slack.File) {
	uploadToImgur(file.URLPrivateDownload, a.SlackAppToken, a.ImgurClientID)
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
					a.handleFile(ev.File)
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
