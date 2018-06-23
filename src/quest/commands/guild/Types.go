package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"bytes"
	"fmt"
)

func Types(session *discordgo.Session, message *discordgo.MessageCreate, _ map[string]string, bot commands.Bot) commands.BotError {
	var buffer bytes.Buffer
	buffer.WriteString("```")
	for k, v := range bot.Regex {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	buffer.WriteString("```")
	session.ChannelMessageSend(message.ChannelID, buffer.String())
	return nil
}
