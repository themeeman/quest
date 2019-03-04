package modlog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
	"time"
)

type CasePurge struct {
	ModeratorID string `json:"moderator_id"`
	ChannelID   string `json:"channel_id"`
	Amount      int64  `json:"amount"`
}

func (cp *CasePurge) Embed(session *discordgo.Session) *discordgo.MessageEmbed {
	moderator := utility.GetUser(session, cp.ModeratorID)
	return &discordgo.MessageEmbed{
		Title:     "Purge Messages",
		Color:     0xffff00,
		Timestamp: utility.TimeToTimestamp(time.Now().UTC()),
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: moderator.AvatarURL(""),
			Name:    moderator.String(),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Channel",
				Value:  "<#" + cp.ChannelID + ">",
				Inline: true,
			},
			{
				Name:   "Amount",
				Value:  fmt.Sprint(cp.Amount),
				Inline: true,
			},
		},
	}
}
