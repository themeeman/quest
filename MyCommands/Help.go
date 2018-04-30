package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"bytes"
	"fmt"
	"strings"
)

func Help(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, context commands.Bot) commands.BotError {
	var data = *context.CommandMap
	if args["Command"] == "" {
		var buf bytes.Buffer
		for name, v := range data {
			buf.WriteString(fmt.Sprintf("**%v - ** %v\n", name, v.Description))
		}
		fields := []*discordgo.MessageEmbedField{
			{
				Name:  "Commands",
				Value: buf.String(),
			},
		}
		session.ChannelMessageSendEmbed(message.ChannelID, context.Embed("Help", "", fields))
	} else {
		cmdName := args["Command"]
		cmdInfo, ok := data[cmdName]
		if !ok {
			return commands.UnknownCommandError{
				Command: cmdName,
			}
		}
		var buffer bytes.Buffer
		buffer.WriteString(context.Prefix + cmdName)
		for _, v := range cmdInfo.Arguments {
			if v.Optional {
				buffer.WriteString(fmt.Sprintf(" <%s>", v.Name))
			} else {
				buffer.WriteString(fmt.Sprintf(" [%s]", v.Name))
			}
		}
		fields := []*discordgo.MessageEmbedField{
			{
				Name:  "Usage",
				Value: "```" + buffer.String() + "```",
			},
		}
		session.ChannelMessageSendEmbed(message.ChannelID, context.Embed(strings.ToTitle(cmdName), cmdInfo.Description, fields))
	}
	return nil
}
