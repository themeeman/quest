package events

import (
	"../modlog"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func (bot BotEvents) GuildCreate(session *discordgo.Session, guild *discordgo.GuildCreate) {
	session.State.GuildAdd(guild.Guild)
	time.Sleep(time.Second)
	g := bot.Guilds.Get(guild.ID)
	fmt.Println(guild.ID)
	if g.Modlog.Valid {
		go modlog.StartLogging(session, g.Modlog, &g.Cases)
	}
}
