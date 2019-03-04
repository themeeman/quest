package events

import (
<<<<<<< HEAD:src/quest/events/ready.go
	commands "../../discordcommands"
	"../db"
	"../modlog"
	"../structures"
	"../utility"
=======
>>>>>>> master:events/ready.go
	"fmt"
	"github.com/bwmarrin/discordgo"
	commands "github.com/tomvanwoow/discordcommands"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/structures"
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
				if member.MuteExpires.Valid {
					member.MuteExpires.Valid = false
					duration := member.MuteExpires.Time.UTC().UnixNano() - now.UnixNano()
					go WaitAndUnmute(session, guild, member, time.Duration(duration))
				}
			}
		}
	}
}

func WaitAndUnmute(session *discordgo.Session, guild *structures.Guild, member *structures.Member, duration time.Duration) {
	time.Sleep(duration)
	m, err := session.GuildMember(guild.ID, member.ID)
	if err != nil {
		return
	}
	fmt.Println(m.Roles)
	if found, _ := utility.Contains(m.Roles, guild.MuteRole.String); found {
		return
	}
	err = session.GuildMemberRoleRemove(guild.ID, member.ID, guild.MuteRole.String)
	if err != nil {
		return
	}
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CaseUnmute{
			ModeratorID: "412702549645328397",
			UserID:      member.ID,
			Reason:      "Auto",
		}
	}
}
