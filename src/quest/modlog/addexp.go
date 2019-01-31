package modlog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type CaseAddExp struct {
	AdminID    string `json:"admin_id"`
	Experience int64  `json:"experience"`
	UserID     string `json:"user_id"`
}

func (c *CaseAddExp) Embed(session *discordgo.Session) *discordgo.MessageEmbed {
	admin := getUser(session, c.AdminID)
	var userName string
	if c.UserID == "" {
		userName = "Themself"
	} else {
		user := getUser(session, c.UserID)
		if user != nil {
			if user.ID == c.AdminID {
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
		Timestamp: timeToTimestamp(time.Now().UTC()),
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
