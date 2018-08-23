package modlog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type CaseBan struct {
	ModeratorID string `json:"moderator_id"`
	UserID      string `json:"user_id"`
	Reason      string `json:"reason"`
}

func (cm *CaseBan) Embed(modlog *Modlog, session *discordgo.Session) *discordgo.MessageEmbed {
	member := getUser(session, cm.UserID)
	moderator := getUser(session, cm.ModeratorID)
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "User",
			Value:  member.String() + " " + member.Mention(),
			Inline: true,
		},
	}
	if cm.Reason != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Reason",
			Value: cm.Reason,
		})
	}
	return &discordgo.MessageEmbed{
		Title:     fmt.Sprintf("Case %d | Ban", len(modlog.Cases)+1),
		Color:     0x880000,
		Timestamp: timeToTimestamp(time.Now().UTC()),
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: moderator.AvatarURL(""),
			Name:    moderator.String(),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: member.AvatarURL(""),
		},
		Fields: fields,
	}
}
