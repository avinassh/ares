package ares

import (
	"log"
	"regexp"
)

type regexResult struct {
	botId  string
	userId string
}

func (a *Ares) performMuteAction(msg string) {
	// check if it is either of mute actions and perform
	// steps accordingly
	if status, userID := isMuteAction(msg, a.BotUserID); status == true {
		if a.isAdminUser(userID) || a.BotUserID == userID {
			// never mute admins and the bot
			return
		}
		a.muteUser(userID)
	} else if status, userID := isUnMuteAction(msg, a.BotUserID); status == true {
		a.unMuteUser(userID)
	}
}

func (a *Ares) muteUser(userID string) {
	log.Println("Muting :", userID)
	a.MutedUsers[userID] = true
}

func (a *Ares) unMuteUser(userID string) {
	log.Println("Unmuting :", userID)
	delete(a.MutedUsers, userID)
}

func (a *Ares) isMuted(userID string) bool {
	_, ok := a.MutedUsers[userID]
	return ok
}

// TODO: Refactor isMuteAction and isUnMuteAction both are almost same

// Check if the message is of the form
// "<@U5WLDJUF0> mute <@U5G04KQ21>" using regex
func isMuteAction(msg, botId string) (bool, string) {
	r, _ := regexp.Compile("<@(\\w{9})> mute <@(\\w{9})>")
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

// Check if the message is of the form
// "<@U5WLDJUF0> mute <@U5G04KQ21>" using regex
func isUnMuteAction(msg, botId string) (bool, string) {
	r, _ := regexp.Compile("<@(\\w{9})> unmute <@(\\w{9})>")
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
