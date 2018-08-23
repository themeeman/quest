package commands

import (
	commands "../../discordcommands"
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sort"
	"strings"
)

func (bot *Bot) Help(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
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
		level := bot.UserGroup(session, guild, author)
		for _, name := range names {
			command := bot.CommandMap[name]
			sufficient := level >= command.Group
			if !command.Hidden && sufficient {
				buf.WriteString(fmt.Sprintf("**%s - ** %s\n", name, command.Description))
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
		cmdInfo, name := getCommand(bot.CommandMap, args["Command"])
		if cmdInfo == nil {
			return UnknownCommandError{
				Command: name,
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

func getCommand(commands commands.CommandMap, name string) (*commands.Command, string) {
	name = strings.ToLower(name)
	command, okc := commands[name]
	if !okc {
		for n, cmd := range commands {
			for _, alias := range cmd.Aliases {
				if name == alias {
					return getCommand(commands, n)
				}
			}
		}
	}
	return command, name
}
