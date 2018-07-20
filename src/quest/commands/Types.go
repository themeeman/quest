package commands

import (
	"github.com/bwmarrin/discordgo"
	"bytes"
	"fmt"
)

func (bot *Bot) Types(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string) error {
	var buffer bytes.Buffer
	buffer.WriteString("```")
	for k, v := range bot.Regex {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	buffer.WriteString("```")
	session.ChannelMessageSend(message.ChannelID, buffer.String())
	return nil
}
