package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/utility"
	"strings"
)

func (bot *Bot) Unmute(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	ch, _ := session.State.Channel(message.ChannelID)
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
	member, _ := session.State.Member(ch.GuildID, user.ID)
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	var found bool
	if !guild.MuteRole.Valid {
		return fmt.Errorf("No mute role has been configured for the server! Use q:set muterole [Value]")
	}
	for _, r := range member.Roles {
		if r == guild.MuteRole.String {
			found = true
		}
	}
	if !found {
		return UnmutedError{
			Username:      user.Username,
			Discriminator: user.Discriminator,
		}
	}
	err := session.GuildMemberRoleRemove(ch.GuildID, user.ID, guild.MuteRole.String)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 403 Forbidden") {
			return fmt.Errorf("Make sure the bot has Manage Roles Permission in Discord!")
		} else {
			return UserNotFoundError{}
		}
	}
	m := guild.Get(user.ID)
	m.MuteExpires.Valid = false
	if args["Reason"] == "" {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully unmuted %s#%s!", user.Username, user.Discriminator), nil))
	} else {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully unmuted %s#%s! Reason: %s", user.Username, user.Discriminator, args["Reason"]), nil))
	}
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CaseUnmute{
			ModeratorID: message.Author.ID,
			UserID:      user.ID,
			Reason:      args["Reason"],
		}
	}
	return nil
}
