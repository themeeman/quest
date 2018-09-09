package modlog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type CasePurge struct {
	ModeratorID string `json:"moderator_id"`
	Amount      int    `json:"amount"`
}

func (cp *CasePurge) Embed(session *discordgo.Session) *discordgo.MessageEmbed {
	moderator := getUser(session, cp.ModeratorID)
	return &discordgo.MessageEmbed{
		Title:     "Purge Messages",
		Color:     0xffff00,
		Timestamp: timeToTimestamp(time.Now().UTC()),
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: moderator.AvatarURL(""),
			Name:    moderator.String(),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Amount",
				Value:  fmt.Sprint(cp.Amount),
				Inline: true,
			},
		},
	}
}
