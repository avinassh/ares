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

func main() {
	api := slack.New(os.Getenv("BOT_TOKEN"))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			if ev.SubType == "file_share" {
				go deleteFile(ev.File.ID)
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
