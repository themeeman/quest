package modlog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type CaseSet struct {
	ModeratorID string
	Option      string
	Value       string
}

func (cs *CaseSet) Embed(modlog *Modlog, session *discordgo.Session) *discordgo.MessageEmbed {
	moderator := getUser(session, cs.ModeratorID)
	return &discordgo.MessageEmbed{
		Title:     fmt.Sprintf("Case %d | Set Option", len(modlog.Cases)+1),
		Color:     0xffffff,
		Timestamp: timeToTimestamp(time.Now().UTC()),
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: moderator.AvatarURL(""),
			Name:    moderator.String(),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Option",
				Value:  cs.Option,
				Inline: true,
			},
			{
				Name:   "Value",
				Value:  cs.Value,
				Inline: true,
			},
		},
	}
}
