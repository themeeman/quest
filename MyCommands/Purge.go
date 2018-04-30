package MyCommands

import (
	"github.com/bwmarrin/discordgo"
	commands "discordcommands"
	"strconv"
	"fmt"
)

func Purge(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string, _ commands.Bot) commands.BotError {
	i, _ := strconv.Atoi(args["Amount"])
	msgs, err := session.ChannelMessages(message.ChannelID, i, "", "", "")
	if err != nil {
		return nil
	}
	ids := make([]string, i)
	for i, v := range msgs {
		ids[i] = v.ID
		fmt.Println(v.Content)
	}
	//session.ChannelMessagesBulkDelete(message.ChannelID, ids)
	return nil
}
