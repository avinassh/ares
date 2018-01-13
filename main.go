package main

import (
	"os"
	"strings"

	"github.com/avinassh/ares/core"
)

func main() {

	a := ares.Ares{
		SlackBotToken: os.Getenv("BOT_TOKEN"),
		SlackAppToken: os.Getenv("APP_TOKEN"),
		ImgurClientID: os.Getenv("IMGUR_CLIENT_ID"),
		SlackAppID:    os.Getenv("APP_ID"),
		Moderators:    strings.Split(os.Getenv("MOD_IDS"), ","),
	}
	a.Run()
}
