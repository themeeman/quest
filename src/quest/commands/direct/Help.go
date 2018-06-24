package direct

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"bytes"
	"sort"
	"fmt"
	"strings"
)

func Help(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot *commands.Bot) commands.BotError {
	if args["Command"] == "" {
		var buf bytes.Buffer
		names := make([]string, len(bot.CommandMap))
		i := 0
		for name := range bot.CommandMap {
			names[i] = name
			i++
		}
		sort.Strings(names)
		for _, name := range names {
			v := bot.CommandMap[name]
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
		if len(cmdInfo.Aliases) > 0 {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Aliases",
				Value: strings.Join(cmdInfo.Aliases, ","),
			})
		}
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed(strings.Title(name), cmdInfo.Description, fields))
	}
	return nil
}
