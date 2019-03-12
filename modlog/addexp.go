package modlog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
	"time"
)

type CaseAddExp struct {
	ModeratorID string `db:"admin_id"`
	Experience  int64  `db:"experience"`
	UserID      string `db:"user_id"`
}

func (c *CaseAddExp) Embed(session *discordgo.Session) *discordgo.MessageEmbed {
	admin := utility.GetUser(session, c.ModeratorID)
	var userName string
	if c.UserID == "" {
		userName = "Themself"
	} else {
		user := utility.GetUser(session, c.UserID)
		if user != nil {
			if user.ID == c.ModeratorID {
				userName = "Themself"
			} else {
				userName = user.String()
			}
		} else {
			userName = fmt.Sprint(`User Not Found (%s)`, user.ID)
		}
	}
	return &discordgo.MessageEmbed{
		Title:     "Add Experience",
		Color:     0x000055,
		Timestamp: utility.TimeToTimestamp(time.Now().UTC()),
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: admin.AvatarURL(""),
			Name:    admin.String(),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User",
				Value:  userName,
				Inline: true,
			},
			{
				Name:   "Experience",
				Value:  fmt.Sprint(c.Experience),
				Inline: true,
			},
		},
	}

}
