package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"strconv"
	"fmt"
	"strings"
)

func Purge(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, _ commands.Bot) commands.BotError {
	i, _ := strconv.Atoi(strings.Replace(args["Amount"], ",", "", -1))
	msgs, err := session.ChannelMessages(message.ChannelID, i+1, "", "", "")
	if err != nil {
		return nil
	}
	ids := make([]string, i+1)
	for i, v := range msgs {
		ids[i] = v.ID
		fmt.Println(v.Content)
	}
	session.ChannelMessagesBulkDelete(message.ChannelID, ids)
	return nil
}
