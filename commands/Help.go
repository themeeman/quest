package commands

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	commands "github.com/tomvanwoow/disgone"
	"github.com/tomvanwoow/quest/utility"
	"sort"
	"strings"
)

func (bot *Bot) Help(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	if args["Command"] == "" {
		var buf bytes.Buffer
		names := make([]string, len(bot.Commands))
		i := 0
		for name := range bot.Commands {
			names[i] = name
			i++
		}
		sort.Strings(names)
		guild, _ := session.Guild(utility.MustGetGuildID(session, message))
		author, _ := session.GuildMember(guild.ID, message.Author.ID)
		level := bot.UserGroup(session, guild, author)
		for _, name := range names {
			command := bot.Commands[name]
			if !command.Hidden && level >= command.Group {
				buf.WriteString(fmt.Sprintf("**%s - ** %s\n", name, command.Description))
			}
		}
		fields := []*discordgo.MessageEmbedField{
			{
				Name:  "Commands",
				Value: buf.String(),
			},
		}
		ch, _ := session.UserChannelCreate(message.Author.ID)
		if ch != nil {
			_, err := session.ChannelMessageSendEmbed(ch.ID, bot.Embed("Help", "", fields))
			if err != nil {
				session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed("Help", "", fields))
			}
		}
	} else {
		cmdInfo, name := getCommand(bot.Commands, args["Command"])
		if cmdInfo == nil {
			return UnknownCommandError{
				Command: name,
			}
		}
		fields := make([]*discordgo.MessageEmbedField, 1, 3)
		fields[0] = &discordgo.MessageEmbedField{
			Name:  "Usage",
			Value: "```" + cmdInfo.GetUsage(bot.Prefix, name) + "```",
		}
		if len(cmdInfo.Examples) > 0 {
			var exampleBuffer bytes.Buffer
			for _, v := range cmdInfo.Examples {
				exampleBuffer.WriteString(fmt.Sprintf("`%s`\n", v))
			}

			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Examples",
				Value:  exampleBuffer.String(),
				Inline: true,
			})
		}
		if len(cmdInfo.Aliases) > 0 {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Aliases",
				Value:  strings.Join(cmdInfo.Aliases, ", "),
				Inline: true,
			})
		}
		session.ChannelMessageSendEmbed(message.ChannelID, bot.Embed(strings.Title(name), cmdInfo.Description, fields))
	}
	return nil
}

func getCommand(commands commands.CommandMap, name string) (*commands.Command, string) {
	name = strings.ToLower(name)
	command, okc := commands[name]
	if okc {
		return command, name
	}
	for n, cmd := range commands {
		for _, alias := range cmd.Aliases {
			if name == alias {
				return getCommand(commands, n)
			}
		}
	}
	return nil, ""
}
