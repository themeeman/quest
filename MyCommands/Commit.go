package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"fmt"
)

func Commit(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string, bot commands.Bot) commands.BotError {
	if message.Author.ID != "164759167561629696" {
		return commands.PermissionsError{}
	}
	err := commands.PostAllGuildData(bot.DB, bot.Guilds)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	session.ChannelMessageSend(message.ChannelID, "Done")
	return nil
}
