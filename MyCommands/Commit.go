package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"fmt"
	"encoding/json"
)

func Commit(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string, bot commands.Bot) commands.BotError {
	if message.Author.ID != "164759167561629696" {
		return commands.PermissionsError{}
	}
	c, _ := session.Channel(message.ChannelID)
	guild, ok := commands.FindGuildByID(bot.Guilds, c.GuildID)
	if ok {
		b, _ := json.MarshalIndent(guild, "", " ")
		session.ChannelMessageSend(message.ChannelID, "```json\n" + string(b) + "```")
	}
	err := commands.PostAllGuildData(bot.DB, bot.Guilds)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	session.ChannelMessageSend(message.ChannelID, "Done")
	return nil
}
