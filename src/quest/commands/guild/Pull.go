package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
)

func Pull(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string, bot *commands.Bot) commands.BotError {
	if message.Author.ID == "164759167561629696" {
		guilds, err := commands.QueryAllData(bot.DB)
		if err == nil {
			bot.Guilds = guilds
			session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
		}
	}
	return nil
}
