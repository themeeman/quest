package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
)

func SetMuteRole(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string, ctx commands.Bot) commands.BotError {
	ch, _ := session.Channel(message.ChannelID)
	_, guild, _, _ := commands.FindGuildByID(ctx.Guilds, ch.GuildID)
	guild.MuteRole.String = message.MentionRoles[0]
	guild.MuteRole.Valid = true
	session.ChannelMessageSend(message.ChannelID, "Success!")
	return nil
}
