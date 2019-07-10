package commands

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
)

func (bot *Bot) Inventory(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	guild.RLock()
	member := guild.Members.Get(message.Author.ID)
	guild.RUnlock()
	var inv bytes.Buffer
	member.RLock()
	for id, q := range member.Chests {
		if q > 0 {
			if chest, ok := bot.Chests[id]; ok {
				inv.WriteString(fmt.Sprintf("%s - %d\n", chest.Name, q))
			}
		}
	}
	member.RUnlock()
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
