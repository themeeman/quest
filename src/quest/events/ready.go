package events

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"fmt"
	commands "../../discordcommands"
	"time"
)

func Ready(bot *commands.Bot) func(*discordgo.Session, *discordgo.Ready) {
	return func(session *discordgo.Session, ready *discordgo.Ready) {
		var err error
		guilds, err := commands.QueryAllData(bot.DB)
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
				err := commands.PostAllData(bot.DB, bot.Guilds)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("Successfully commited all data")
				}
			}
		}()
		applyMutes(bot, session)
	}
}

func applyMutes(bot *commands.Bot, session *discordgo.Session) {
	now := time.Now().UTC()
	for _, guild := range bot.Guilds {
		if guild.MuteRole.Valid {
			for _, member := range guild.Members {
				if member.MuteExpires.Valid && member.MuteExpires.Time.After(now) {
					go func(guild *commands.Guild, member *commands.Member) {
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