package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"bytes"
	"fmt"
	"strings"
	"encoding/json"
)

func Help(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot commands.Bot) commands.BotError {
	if args["Command"] == "" {
		var buf bytes.Buffer
		for name, v := range bot.CommandMap {
			if _, ok := bot.HandlerMap[name]; ok && !v.Hidden {
				buf.WriteString(fmt.Sprintf("**%s - ** %s\n", name, v.Description))
			}
		}
		fields := []*discordgo.MessageEmbedField{
			{
				Name:  "Commands",
				Value: buf.String(),
			},
		}
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Help", "", fields))
	} else {
		cmdInfo, hand, name := commands.GetCommand(bot, args["Command"])
		if cmdInfo == nil || hand == nil {
			return commands.UnknownCommandError{
				Command: name,
			}
		}
		if message.Author.ID == "164759167561629696" {
			b, err := json.MarshalIndent(cmdInfo, "", " ")
			if err != nil {
				fmt.Println(err)
			}
			session.ChannelMessageSend(message.ChannelID, "```json\n"+string(b)+"```")
		}
		var buffer bytes.Buffer
		buffer.WriteString(bot.Prefix + name)
		for _, v := range cmdInfo.Arguments {
			if v.Optional {
				buffer.WriteString(fmt.Sprintf(" <%s>", v.Name))
			} else {
				buffer.WriteString(fmt.Sprintf(" [%s]", v.Name))
			}
		}
		var fields []*discordgo.MessageEmbedField
		if len(cmdInfo.Examples) > 0 {
			var exampleBuffer bytes.Buffer
			for _, v := range cmdInfo.Examples {
				exampleBuffer.WriteString(fmt.Sprintf("`%s`\n", v))
			}

			fields = []*discordgo.MessageEmbedField{
				{
					Name:  "Usage",
					Value: "```" + buffer.String() + "```",
				},
				{
					Name:  "Examples",
					Value: exampleBuffer.String(),
				},
			}
		} else {
			fields = []*discordgo.MessageEmbedField{
				{
					Name:  "Usage",
					Value: "```" + buffer.String() + "```",
				},
			}
		}
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed(strings.Title(name), cmdInfo.Description, fields))
	}
	return nil
}
