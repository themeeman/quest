package events

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"fmt"
	"../structures"
	"../db"
	"time"
	quest "../commands"
	commands "../../discordcommands"
)

func Ready(bot *quest.Bot) func(*discordgo.Session, *discordgo.Ready) {
	return func(session *discordgo.Session, event *discordgo.Ready) {
		if bot.ReadyEvent {
			return
		}
		bot.ReadyEvent = true
		var err error
		guilds, err := db.QueryAllData(bot.DB)
		if err != nil {
			log.Println("b", err)
		}
		for _, v := range guilds {
			fmt.Println(v)
		}
		bot.Guilds = guilds
		session.UpdateStatus(0, "q:help")
		go func() {
			for {
				time.Sleep(time.Minute * 10)
				err := db.PostAllData(bot.DB, bot.Guilds)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("Successfully commited all data")
				}
			}
		}()
		go func() {
			for {
				err := <-bot.Errors
				if err.Err != nil {
					if e, ok := err.Err.(commands.ZeroArgumentsError); ok {
						bot.Help(session, err.MessageCreate, map[string]string{"Command": e.Command})
					} else {
						session.ChannelMessageSendEmbed(err.ChannelID, commands.ErrorEmbed(err.Err))
					}
				}
			}
		}()
		applyMutes(guilds, session)
	}
}

func applyMutes(guilds structures.Guilds, session *discordgo.Session) {
	now := time.Now().UTC()
	for _, guild := range guilds {
		if guild.MuteRole.Valid {
			for _, member := range guild.Members {
				if member.MuteExpires.Valid && member.MuteExpires.Time.After(now) {
					go func(guild *structures.Guild, member *structures.Member) {
						dur := member.MuteExpires.Time.UTC().UnixNano() - now.UnixNano()
						time.Sleep(time.Duration(dur))
						session.GuildMemberRoleRemove(guild.ID, member.ID, guild.MuteRole.String)
					}(guild, member)
				} else if member.MuteExpires.Valid && member.MuteExpires.Time.Before(now) {
					go session.GuildMemberRoleRemove(guild.ID, member.ID, guild.MuteRole.String)
				}
			}
		}
	}
}
