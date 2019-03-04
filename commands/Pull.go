package commands

import (
	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) Pull(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string) error {
	if message.Author.ID == "164759167561629696" {
		guilds, err := db.QueryAllData(bot.DB)
		if err == nil {
			bot.Guilds = guilds
			session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
		}
	}
	return nil
}
