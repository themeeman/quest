package events

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"fmt"
	"strings"
	"time"
	"runtime/debug"
	commands "../../discordcommands"
	"../experience"
)

func MessageCreate(bot *commands.Bot) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(session *discordgo.Session, message *discordgo.MessageCreate) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(string(debug.Stack()))
				session.ChannelMessageSend(message.ChannelID, "```"+ `An unexpected panic occured in the execution of that command.
`+ fmt.Sprint(r) + "\nTry again later, or contact themeeman#8354" + "```")
			}
		}()
		fmt.Println(message.Content)
		if !message.Author.Bot {
			if strings.ToLower(message.Content) == "good bot" {
				m, _ := session.ChannelMessageSend(message.ChannelID, "Your compliments mean nothing to me")
				time.Sleep(5 * time.Second)
				if m != nil {
					session.ChannelMessageDelete(message.ChannelID, m.ID)
				}
			}
			if commands.HasPrefix(message.Content, bot.Prefix) {
				err := commands.ExecuteCommand(session, message, bot)
				if err != nil {
					session.ChannelMessageSendEmbed(message.ChannelID, commands.ErrorEmbed(err))
				}
			}
			experience.GrantExp(bot, session, message)
		}
	}
}