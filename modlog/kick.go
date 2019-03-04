package modlog

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
	"time"
)

type CaseKick struct {
	ModeratorID string `json:"moderator_id"`
	UserID      string `json:"user_id"`
	Reason      string `json:"reason"`
}

func (cm *CaseKick) Embed(session *discordgo.Session) *discordgo.MessageEmbed {
	member := utility.GetUser(session, cm.UserID)
	moderator := utility.GetUser(session, cm.ModeratorID)
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
		Title:     "Kick",
		Color:     0xffcccc,
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
