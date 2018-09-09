package modlog

import "github.com/bwmarrin/discordgo"

type CaseAddExp struct {
}

func (c *CaseAddExp) Embed(modlog *Modlog, session *discordgo.Session) *discordgo.MessageEmbed {
}