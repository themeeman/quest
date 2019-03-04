package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/utility"
)

func (bot *Bot) Warn(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	var user *discordgo.User
	if len(args["User"]) == 18 {
		var err error
		user, err = session.User(args["User"])
		if err != nil {
			return UserNotFoundError{}
		}
	} else if len(message.Mentions) > 0 {
		user = message.Mentions[0]
	} else {
		return UserNotFoundError{}
	}
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CaseWarn{
			ModeratorID: message.Author.ID,
			UserID:      user.ID,
			Reason:      args["Reason"],
		}
	}
	g, err := session.Guild(guild.ID)
	if err != nil {
		return nil
	}
	ch, err := session.UserChannelCreate(user.ID)
	if err == nil {
		if args["Reason"] == "" {
			session.ChannelMessageSend(ch.ID, fmt.Sprintf("You were warned in **%s**", g.Name))
		} else {
			session.ChannelMessageSend(ch.ID, fmt.Sprintf("You were warned in **%s** for reason: %s", g.Name, args["Reason"]))
		}
	}
	return nil
}
