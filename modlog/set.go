package modlog

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
	"time"
)

type CaseSet struct {
	ModeratorID string `json:"moderator_id"`
	Option      string `json:"option"`
	Value       string `json:"value"`
}

func (cs *CaseSet) Embed(session *discordgo.Session) *discordgo.MessageEmbed {
	moderator := utility.GetUser(session, cs.ModeratorID)
	return &discordgo.MessageEmbed{
		Title:     "Set Option",
		Color:     0xffffff,
		Timestamp: utility.TimeToTimestamp(time.Now().UTC()),
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
