package ares

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

// I don't like templating in Go
var msgTemplate string = "Hi %s(%s), welcome to the Dev Up (Devs and Hackers) Slack team!\n\nPlease introduce yourself in the <#C0CMVBB2Q|intro> channel. Tell us about your work, your open source projects(if any), the languages/tools/frameworks you use.\n\nWe have more than 50 channels - feel free to explore them in the sidebar. Some of the more popular ones include <#C0D4RA5HV|python>, <#C0D4QK537|js>, <#C0JP0FDMH|androiddev>, <#C16R4U2D6|swift>, <#C27K2UWH2|careers>, <#C0FQD6KG8|showcase-projects>.\n\nWe also have non-tech channels like <#C0G273YCT|books>, <#C164BNBRP|television>, <#C0U1W6NKT|health>, <#C0CNQCE73|business>. Also, city specific channels like <#C0F9Z6J85|bangalore>, <#C0JKT0P4N|hyderabad>, <#C0G7RTC5N|delhi-ncr> where we plan local events and meetups. We have private channels like #television-spoiled, #sports, #english, #boardgames, #wisdom among others. To join these channels, drop a message in <#C06VBR8UT|general> and someone will add you.\n\nWhile we do not have any strict rules, please do pay attention to our guidelines over here: http://bit.ly/devs-rule.\n\nFeel free to contact <@U06VBQ8TB|v> in case you need any help or have any suggestions.\n\nEnjoy your stay! Keep Hacking :nerd_face:"

func (a *Ares) onBoardUser(newUser slack.User) {
	api := slack.New(a.SlackBotToken)
	msg := fmt.Sprintf(msgTemplate, newUser.RealName, newUser.Name)
	params := slack.PostMessageParameters{}
	if _, _, err := api.PostMessage(newUser.ID, msg, params); err != nil {
		log.Println("Failed to send welcome message to new user: ", newUser.Name)
		return
	}
	a.notifyAdmin(newUser)
}

func (a *Ares) notifyAdmin(newUser slack.User) {
	api := slack.New(a.SlackBotToken)
	msg := fmt.Sprintf("User: %s successfully on-boarded", newUser.Name)
	params := slack.PostMessageParameters{}
	for _, admin := range a.Admins {
		if _, _, err := api.PostMessage(admin, msg, params); err != nil {
			log.Println("Failed to notify admin about new user: ", newUser.Name)
			return
		}
	}
}
