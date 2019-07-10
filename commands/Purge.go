package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/utility"
	"strconv"
	"strings"
	"time"
)

func (bot *Bot) Purge(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	i, _ := strconv.Atoi(strings.Replace(args["Amount"], ",", "", -1))
	if i >= 200 {
		return fmt.Errorf("I can't purge more than 200 messages")
	}
	if i == 0 {
		return fmt.Errorf("Invalid purge value")
	}
	msgs, err := session.ChannelMessages(message.ChannelID, i, message.ID, "", "")
	if err != nil {
		return nil
	}
	fmt.Println(msgs)
	ids := make([]string, len(msgs))
	for i, v := range msgs {
		ids[i] = v.ID
		fmt.Println(v.Content)
	}
	err = session.ChannelMessagesBulkDelete(message.ChannelID, ids)
	if err != nil {
		return fmt.Errorf("It seems I do not have permissions to delete messages!")
	} else {
		session.ChannelMessageDelete(message.ChannelID, message.ID)
		m, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("☑ Successfully purged %d messages", len(msgs)))
		if err == nil {
			go func() {
				time.Sleep(time.Second * 10)
				session.ChannelMessageDelete(m.ChannelID, m.ID)
			}()
		}
	}
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	guild.RLock()
	defer guild.RUnlock()
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CasePurge{
			ModeratorID: message.Author.ID,
			ChannelID:   message.ChannelID,
			Amount:      len(msgs),
		}
	}
	return nil
}
