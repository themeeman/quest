package commands

import (
	commands "../../discordcommands"
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) Inventory(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	member := guild.Members.Get(message.Author.ID)
	var inv bytes.Buffer
	for id, q := range member.Chests {
		if q > 0 {
			if chest, ok := bot.Chests[id]; ok {
				inv.WriteString(fmt.Sprintf("%s - %d", chest.Name, q))
			}
		}
	}
	if inv.String() != "" {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("", "", []*discordgo.MessageEmbedField{
			{
				Name:  "Inventory",
				Value: inv.String(),
			},
		}))
	} else {
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Sorry!", "Sorry buddy, you don't have any items!", nil))
	}
	return nil
}
