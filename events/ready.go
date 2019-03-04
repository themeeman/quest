package events

import (
	commands "../../discordcommands"
	"../db"
	"../modlog"
	"../structures"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func (bot BotEvents) Ready(session *discordgo.Session, _ *discordgo.Ready) {
	var err error
	guilds, err := db.QueryAllData(bot.DB)
	if err != nil {
		log.Fatalln("b", err)
	}
	for _, v := range guilds {
		//fmt.Println(v)
		fmt.Println(v.ID, v.Cases)
	}
	if guilds == nil {
		guilds = make(structures.Guilds)
	}
	bot.Guilds = guilds
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
		var err struct {
			Err error
			*discordgo.MessageCreate
		}
		for {
			err = <-bot.Errors
			if err.Err != nil {
				fmt.Println(err)
				if e, ok := err.Err.(commands.ZeroArgumentsError); ok {
					bot.Help(session, err.MessageCreate, map[string]string{"Command": e.Command})
				} else {
					session.ChannelMessageSendEmbed(err.ChannelID, bot.ErrorEmbed(err.Err))
				}
			}
		}
	}()
	go applyMutes(guilds, session)
}

func applyMutes(guilds structures.Guilds, session *discordgo.Session) {
	now := time.Now().UTC()
	for _, guild := range guilds {
		if guild.MuteRole.Valid {
			for _, member := range guild.Members {
				if member.MuteExpires.Valid && member.MuteExpires.Time.After(now) {
					go func(guild *structures.Guild, member *structures.Member) {
						member.MuteExpires.Valid = false
						dur := member.MuteExpires.Time.UTC().UnixNano() - now.UnixNano()
						time.Sleep(time.Duration(dur))
						err := session.GuildMemberRoleRemove(guild.ID, member.ID, guild.MuteRole.String)
						if err == nil && guild.Modlog.Valid {
							guild.Modlog.Log <- &modlog.CaseUnmute{
								ModeratorID: "412702549645328397",
								UserID:      member.ID,
								Reason:      "Auto",
							}
						}
					}(guild, member)
				} else if member.MuteExpires.Valid && member.MuteExpires.Time.Before(now) {
					go func(guild *structures.Guild, member *structures.Member) {
						member.MuteExpires.Valid = false
						err := session.GuildMemberRoleRemove(guild.ID, member.ID, guild.MuteRole.String)
						if err == nil && guild.Modlog.Valid {
							guild.Modlog.Log <- &modlog.CaseUnmute{
								ModeratorID: "412702549645328397",
								UserID:      member.ID,
								Reason:      "Auto",
							}
						}
					}(guild, member)
				}
			}
		}
	}
}
