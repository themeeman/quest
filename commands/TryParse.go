package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
)

func (bot *Bot) TryParse(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	pattern, ok := bot.Types[args["Type"]]
	if !ok {
		return fmt.Errorf(`The provided argument for the Type was incorrect:
%s is **not** a Type.
Use q:types to view all types.`, args["Type"])
	}
	result, err := regexp.MatchString(pattern, args["Value"])
	if err != nil {
		return nil
	}
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Type",
			Value:  args["Type"],
			Inline: true,
		},
		{
			Name:   "Value",
			Value:  args["Value"],
			Inline: true,
		},
		{
			Name:  "Result",
			Value: fmt.Sprintf("%t", result),
		},
	}
	session.ChannelMessageSendEmbed(message.ChannelID,
		bot.Embed("Result", "", fields))
	return nil
}
