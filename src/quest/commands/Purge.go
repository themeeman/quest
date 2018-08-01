package commands

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"fmt"
	"strings"
	"time"
)

func (bot *Bot) Purge(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	i, _ := strconv.Atoi(strings.Replace(args["Amount"], ",", "", -1))
	if i >= 200 {
		return fmt.Errorf("I can't purge more than 200 messages")
	}
	msgs, err := session.ChannelMessages(message.ChannelID, i, message.ID, "", "")
	if err != nil {
		return nil
	}
	ids := make([]string, i)
	for i, v := range msgs {
		ids[i] = v.ID
		fmt.Println(v.Content)
	}
	err = session.ChannelMessagesBulkDelete(message.ChannelID, ids)
	if err != nil {
		return fmt.Errorf("It seems I do not have permissions to delete messages!")
	} else {
		session.ChannelMessageDelete(message.ChannelID, message.ID)
		m, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("☑ Successfully purged %d messages", i))
		if err == nil {
			go func() {
				time.Sleep(time.Second * 10)
				session.ChannelMessageDelete(m.ChannelID, m.ID)
			}()
		}
	}
	return nil
}
