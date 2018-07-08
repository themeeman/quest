package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"bytes"
	"fmt"
	"strings"
	"encoding/json"
	"sort"
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
		guild, _ := session.Guild(commands.MustGetGuildID(session, message))
		author, _ := session.GuildMember(guild.ID, message.Author.ID)
		for _, name := range names {
			v := bot.CommandMap[name]
			sufficient, _, _ := commands.SufficentPermissions(session, guild, author, bot, v)
			if _, ok := bot.HandlerMap[name]; ok && !v.Hidden && sufficient {
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
		var fields []*discordgo.MessageEmbedField
		if len(cmdInfo.Examples) > 0 {
			var exampleBuffer bytes.Buffer
			for _, v := range cmdInfo.Examples {
				exampleBuffer.WriteString(fmt.Sprintf("`%s`\n", v))
			}

			fields = []*discordgo.MessageEmbedField{
				{
					Name:  "Usage",
					Value: "```" + cmdInfo.GetUsage(bot.Prefix, name) + "```",
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
					Value: "```" + cmdInfo.GetUsage(bot.Prefix, name) + "```",
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
