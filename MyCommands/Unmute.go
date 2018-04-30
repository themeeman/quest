package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"strings"
	"fmt"
)

func Unmute(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, context commands.Bot) commands.BotError {
	ch, _ := session.State.Channel(message.ChannelID)
	var user *discordgo.User
	if len(message.Mentions) == 0 {
		user = new(discordgo.User)
	} else {
		user = message.Mentions[0]
	}
	member, _ := session.State.Member(ch.GuildID, user.ID)
	var found bool
	for _, r := range member.Roles {
		if r == "413273250131345409" {
			found = true
		}
	}
	if !found {
		return commands.UnmutedError{
			Username:      user.Username,
			Discriminator: user.Discriminator,
		}
	}
	err := session.GuildMemberRoleRemove(ch.GuildID, user.ID, "413273250131345409")
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 403 Forbidden") {
			return commands.PermissionsError{}
		} else {
			return commands.UserNotFoundError{}
		}
	}
	if args["Reason"] == "" {
		session.ChannelMessageSendEmbed(message.ChannelID, context.Embed("Success!", fmt.Sprintf("Successfully unmuted %session#%session!", user.Username, user.Discriminator), nil))
	} else {
		session.ChannelMessageSendEmbed(message.ChannelID, context.Embed("Success!", fmt.Sprintf("Successfully unmuted %session#%session! Reason: %session", user.Username, user.Discriminator, args["Reason"]), nil))
	}
	return nil
}
