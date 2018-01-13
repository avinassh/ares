package ares

import (
	"log"
	"regexp"

	"github.com/nlopes/slack"
)

func (a *Ares) performKickAction(msg, channelId string) {
	// check if it is either of mute actions and perform
	// steps accordingly
	if status, userID := isKickAction(msg, a.BotUserID); status == true {
		if a.isAdminUser(userID) {
			// never kick admins
			return
		}
		a.kickUser(userID, channelId)
	}
}

func (a *Ares) kickUser(userID, channelId string) {
	log.Printf("Kicking : %s from %s", userID, channelId)
	api := slack.New(a.SlackAppToken)
	if err := api.KickUserFromGroup(channelId, userID); err != nil {
		log.Println("failed to kick the user: ", err)
	}
}

// Check if the message is of the form
// "<@U5WLDJUF0> mute <@U5G04KQ21>" using regex
func isKickAction(msg, botId string) (bool, string) {
	r, _ := regexp.Compile("<@(\\w{9})> kick <@(\\w{9})>")
	var isMatch bool
	var user string
	if r.MatchString(msg) {
		result := r.FindStringSubmatch(msg)
		// check if the message was indeed sent to bot
		if result[1] == botId {
			user = result[2]
			isMatch = true
		}

	}
	return isMatch, user
}
