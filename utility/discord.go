package utility

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func MustGetGuildID(session *discordgo.Session, message *discordgo.MessageCreate) string {
	c, err := session.Channel(message.ChannelID)
	if err != nil {
		return ""
	}
	return c.GuildID
}

func TimeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

func GetRole(session *discordgo.Session, guildID string, id string) (*discordgo.Role, error) {
	rs, err := session.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, fmt.Errorf("role not found %s", id)
}


func GetUser(session *discordgo.Session, id string) *discordgo.User {
	m, _ := session.User(id)
	return m
}
