package main

import (
	"os"

	"github.com/avinassh/ares/ares"
)

func main() {
	a := ares.Ares{
		SlackBotToken: os.Getenv("BOT_TOKEN"),
		SlackAppToken: os.Getenv("APP_TOKEN"),
		ImgurClientID: os.Getenv("IMGUR_CLIENT_ID"),
	}
	a.Run()
}
