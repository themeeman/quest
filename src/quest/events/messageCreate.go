package events

import (
	"../experience"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func (bot BotEvents) MessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Println(message.Content)
	if !message.Author.Bot {
		if strings.ToLower(message.Content) == "good bot" {
			m, _ := session.ChannelMessageSend(message.ChannelID, "Your compliments mean nothing to me")
			time.Sleep(5 * time.Second)
			if m != nil {
				session.ChannelMessageDelete(message.ChannelID, m.ID)
			}
		}
		experience.GrantExp(bot.Bot, session, message)
	}
}
