package events

import "github.com/bwmarrin/discordgo"

func GuildCreate(session *discordgo.Session, guild *discordgo.GuildCreate) {
	session.State.GuildAdd(guild.Guild)
}
