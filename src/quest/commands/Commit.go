package commands

import (
	"github.com/bwmarrin/discordgo"
	"../db"
	"fmt"
)

func (bot *Bot) Commit(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string) error {
	if message.Author.ID != "164759167561629696" {
		return nil
	}
	err := db.PostAllData(bot.DB, bot.Guilds)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
