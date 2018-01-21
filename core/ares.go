package ares

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

type Ares struct {
	SlackAppToken string
	SlackBotToken string
	SlackAppID    string
	ImgurClientID string
	BotUserID     string
	// Bot won't be added or re-added to these channels
	ManagedChannels []string
	// list of userIDs of admins
	Admins []string
	// list of userIDs of moderators
	Moderators []string
	// maintain a user dict
	Users map[string]string
	// Muted users list
	MutedUsers map[string]bool
}

func (a *Ares) initBot() {
	api := slack.New(a.SlackAppToken)
	a.getBotandAdmin()

	channels, err := api.GetChannels(true)

	if err != nil {
		log.Fatal("Failed to get public channels list: ", err.Error())
	}

	for _, channel := range channels {
		a.addBotChannel(channel.ID)
	}

	groups, err := api.GetGroups(true)

	if err != nil {
		log.Fatal("Failed to get private channels list: ", err.Error())
	}

	for _, group := range groups {
		a.addBotGroup(group.ID)
	}
}

// Fetches the bot user ID and also saves the admin user ids
func (a *Ares) getBotandAdmin() {
	api := slack.New(a.SlackAppToken)
	users, err := api.GetUsers()

	if err != nil {
		log.Fatal("Failed to fetch slack users info", err.Error())
	}

	a.Users = make(map[string]string)
	a.MutedUsers = make(map[string]bool)

	for _, user := range users {

		if user.Profile.ApiAppID == a.SlackAppID {
			a.BotUserID = user.ID
		}

		if user.IsAdmin {
			a.Admins = append(a.Admins, user.ID)
		} else {
			a.Users[user.Name] = user.ID
		}
	}

	if a.BotUserID == "" {
		log.Fatal("Unable to find bot user on the Slack")
	}
}

func (a *Ares) deleteMsg(channel, msgTimestamp string) {
	api := slack.New(a.SlackAppToken)

	if _, _, err := api.DeleteMessage(channel, msgTimestamp); err != nil {
		log.Println("Failed to delete msg: ", err.Error())
	}
}

func (a *Ares) deleteFile(fileId string) {
	api := slack.New(a.SlackAppToken)

	if err := api.DeleteFile(fileId); err != nil {
		fmt.Println("Deletion of the file failed:", err.Error())
	}
}

func (a *Ares) addBotChannel(channelID string) {
	api := slack.New(a.SlackAppToken)
	if _, err := api.InviteUserToChannel(channelID, a.BotUserID); err != nil {
		if err.Error() != "already_in_channel" {
			log.Println(fmt.Sprintf("Failed to add bot to %s: %s", channelID, err.Error()))
		}
	}
}

func (a *Ares) addBotGroup(group string) {
	api := slack.New(a.SlackAppToken)
	if _, _, err := api.InviteUserToGroup(group, a.BotUserID); err != nil {
		log.Println(fmt.Sprintf("Failed to add bot to %s: %s", group, err.Error()))
	}
}

func (a *Ares) notifyUser(user, deleteLink string) {
	api := slack.New(a.SlackBotToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color:   "danger",
		Pretext: "Hey there, image uploads to this slack is disabled. Please use external image upload services next time. This incident has been reported :squirrel:",
		Text:    "I have uploaded the image to Imgur for now.",
		Fields: []slack.AttachmentField{
			{
				Title: "If you wish to delete the image uploaded, click on the following link",
				Value: deleteLink,
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(user, "", params)
	if err != nil {
		log.Println("Failed to send DM to user", err.Error())
	}
}

func (a *Ares) sendImgToSlack(channel, user, url, commentText string) {
	api := slack.New(a.SlackBotToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color:    "#D3D3D3",
		Text:     fmt.Sprintf("image originally uploaded by <@%s>", user),
		ImageURL: url,
		Fields: []slack.AttachmentField{
			{
				Title: "",
				Value: commentText,
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(channel, "", params)
	if err != nil {
		log.Println("Failed to send DM to user", err.Error())
	}
}

func (a *Ares) handleFile(file *slack.File, channel string) {
	commentText := file.InitialComment.Comment

	resp := uploadToImgur(file.URLPrivateDownload, a.SlackAppToken, a.ImgurClientID)
	if resp.Status == false {
		log.Println("Failed to download/upload: ", file.ID)
		log.Println("Download url: ", file.URLPrivateDownload)
		return
	}

	a.notifyUser(file.User, resp.DeleteLink)
	a.sendImgToSlack(channel, file.User, resp.Link, commentText)
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

	// run clean up duties
	go a.ClearImages()

	a.initBot()
	log.Println("Bot initialized. Starting moderation duty.")

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:

			if a.isAdminUser(ev.Msg.User) || a.isModUser(ev.Msg.User) {
				a.performMuteAction(ev.Msg.Text)
				a.performKickAction(ev.Msg.Text, ev.Msg.Channel)
			} else if a.isMuted(ev.Msg.User) {
				go a.deleteMsg(ev.Msg.Channel, ev.Msg.Timestamp)
			}

			if ev.SubType == "file_share" {
				if isImageFile(ev.File.Filetype) {
					// TODO: Use worker pool
					go a.handleFile(ev.File, ev.Channel)
				}
			}

		case *slack.TeamJoinEvent:
			if !ev.User.IsBot {
				a.onBoardUser(ev.User)
			}

		case *slack.GroupLeftEvent:
			log.Println("Bot was removed from private channel: ", ev.Channel)
			a.addBotGroup(ev.Channel)

		case *slack.ChannelLeftEvent:
			log.Println("Bot was removed from public channel: ", ev.Channel)
			a.addBotChannel(ev.Channel)

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
