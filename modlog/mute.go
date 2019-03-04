package modlog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
	"time"
)

type CaseMute struct {
	ModeratorID string `json:"moderator_id"`
	UserID      string `json:"user_id"`
	Duration    int64  `json:"duration"`
	Reason      string `json:"reason"`
}

func (cm *CaseMute) Embed(session *discordgo.Session) *discordgo.MessageEmbed {
	member := utility.GetUser(session, cm.UserID)
	moderator := utility.GetUser(session, cm.ModeratorID)
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "User",
			Value:  member.String() + " " + member.Mention(),
			Inline: true,
		},
		{
			Name:   "Duration",
			Value:  fmt.Sprintf("%d Minutes", cm.Duration),
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
		Title:     "Mute",
		Color:     0x00ccff,
		Timestamp: utility.TimeToTimestamp(time.Now().UTC()),
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
