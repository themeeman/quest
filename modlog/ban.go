package modlog

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
	"time"
)

type CaseBan struct {
	ModeratorID string `db:"moderator_id"`
	UserID      string `db:"user_id"`
	Reason      string `db:"reason"`
}

func (cm *CaseBan) Embed(session *discordgo.Session) *discordgo.MessageEmbed {
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
		Title:     "Ban",
		Color:     0x880000,
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
