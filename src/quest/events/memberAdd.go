package events

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

func (bot BotEvents) MemberAdd(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	guild := bot.Guilds.Get(event.GuildID)
	if guild.Autorole.Valid {
		session.GuildMemberRoleAdd(event.GuildID, event.Member.User.ID, guild.Autorole.String)
	}
	if guild.MuteRole.Valid {
		member := guild.Members.Get(event.Member.User.ID)
		if member.MuteExpires.Valid && member.MuteExpires.Time.UTC().After(time.Now().UTC()) {
			err := session.GuildMemberRoleAdd(guild.ID, member.ID, guild.MuteRole.String)
			if err != nil {
				go func() {
					time.Sleep(member.MuteExpires.Time.UTC().Sub(time.Now().UTC()))
					session.GuildMemberRoleRemove(guild.ID, member.ID, guild.MuteRole.String)
				}()
			}
		}
	}
}
