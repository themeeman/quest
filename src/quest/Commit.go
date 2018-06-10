package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"fmt"
)

func Commit(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string, bot commands.Bot) commands.BotError {
	if message.Author.ID != "164759167561629696" {
		return nil
	}
	err := commands.PostAllData(bot.DB, bot.Guilds)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
