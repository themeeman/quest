package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands ".././discordcommands"
)

func SetMuteRole(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	ch, _ := session.Channel(message.ChannelID)
	guild := bot.Guilds.Get(ch.GuildID)
	var role string
	if len(message.MentionRoles) > 0 {
		role = message.MentionRoles[0]
	} else {
		role = args["Role"]
	}
	st, err := session.GuildRoles(commands.MustGetGuildID(session, message))
	if err != nil {
		return commands.UnknownError{}
	}
	var found bool
	for _, v := range st {
		if v.ID == role {
			found = true
		}
	}
	if !found {
		return commands.RoleError{}
	}
	guild.MuteRole.String = role
	guild.MuteRole.Valid = true
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}