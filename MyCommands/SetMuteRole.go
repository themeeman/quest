package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
)

func SetMuteRole(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string, bot commands.Bot) commands.BotError {
	ch, _ := session.Channel(message.ChannelID)
	guild, _ := commands.FindGuildByID(bot.Guilds, ch.GuildID)
	guild.MuteRole.String = message.MentionRoles[0]
	guild.MuteRole.Valid = true
	session.ChannelMessageSend(message.ChannelID, "Success!")
	return nil
}
