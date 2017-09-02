package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

func deleteFile(fileId string) {
	api := slack.New(os.Getenv("APP_TOKEN"))
	if err := api.DeleteFile(fileId); err != nil {
		fmt.Println("Deletion of the file failed:", err.Error())
	}
}

func handleFile(file *slack.File) {
	uploadToImgur(file.URLPrivateDownload, os.Getenv("APP_TOKEN"))
	deleteFile(file.ID)
}

func isImageFile(fileType string) bool {
	for _, fType := range []string{"jpg", "jpeg", "png", "gif"} {
		if fType == fileType {
			return true
		}
	}
	return false
}

func main() {
	api := slack.New(os.Getenv("BOT_TOKEN"))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			if ev.SubType == "file_share" {
				if isImageFile(ev.File.Filetype) {
					handleFile(ev.File)
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
