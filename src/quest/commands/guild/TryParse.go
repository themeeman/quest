package guild

import (
	"github.com/bwmarrin/discordgo"
	commands "../../../discordcommands"
	"regexp"
	"fmt"
)

func TryParse(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, bot *commands.Bot) commands.BotError {
	pattern, ok := bot.Regex[args["Type"]]
	if !ok {
		return commands.TypeError{Name: args["Type"]}
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
