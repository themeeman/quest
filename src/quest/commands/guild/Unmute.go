package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"strings"
	"fmt"
)

func Unmute(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	ch, _ := session.State.Channel(message.ChannelID)
	var user *discordgo.User
	if len(args["User"]) == 18 {
		var err error
		user, err = session.User(args["User"])
		if err != nil {
			return commands.UserNotFoundError{}
		}
	} else if len(message.Mentions) > 0 {
		user = message.Mentions[0]
	} else {
		return commands.UserNotFoundError{}
	}
	member, _ := session.State.Member(ch.GuildID, user.ID)
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	var found bool
	if !guild.MuteRole.Valid {
		return commands.MuteRoleError{}
	}
	for _, r := range member.Roles {
		if r == guild.MuteRole.String {
			found = true
		}
	}
	if !found {
		return commands.UnmutedError{
			Username:      user.Username,
			Discriminator: user.Discriminator,
		}
	}
	err := session.GuildMemberRoleRemove(ch.GuildID, user.ID, guild.MuteRole.String)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 403 Forbidden") {
			return commands.BotPermissionsError{}
		} else {
			return commands.UserNotFoundError{}
		}
	}
	m := guild.Get(user.ID)
	m.MuteExpires.Valid = false
	if args["Reason"] == "" {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully unmuted %session#%session!", user.Username, user.Discriminator), nil))
	} else {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Success!", fmt.Sprintf("Successfully unmuted %session#%session! Reason: %session", user.Username, user.Discriminator, args["Reason"]), nil))
	}
	return nil
}
