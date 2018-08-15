package structures

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"time"
	"database/sql/driver"
	"database/sql"
	"sync"
)

type Modlog struct {
	ChannelID string
	Valid     bool
	Log       chan Case
	Cases     []Case
	Mutex     *sync.Mutex
}

func (m *Modlog) Scan(value interface{}) error {
	null := sql.NullString{}
	err := null.Scan(value)
	if err != nil {
		return err
	}
	m.ChannelID = null.String
	m.Valid = null.Valid
	if m.Valid {
		m.Cases = make([]Case, 0)
		m.Log = make(chan Case)
		m.Mutex = &sync.Mutex{}
	}
	return nil
}

func (m Modlog) Value() (driver.Value, error) {
	null := sql.NullString{
		String: m.ChannelID,
		Valid:  m.Valid,
	}
	return null.Value()
}

type Case interface {
	Embed(Modlog, *discordgo.Session) *discordgo.MessageEmbed
}

func getUser(session *discordgo.Session, id string) (user *discordgo.User) {
	m, _ := session.User(id)
	return m
}

func timeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

type CaseMute struct {
	ModeratorID string
	UserID      string
	Duration    int
	Reason      string
}

func (cm *CaseMute) Embed(modlog Modlog, session *discordgo.Session) *discordgo.MessageEmbed {
	member := getUser(session, cm.UserID)
	moderator := getUser(session, cm.ModeratorID)
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
		Title:     fmt.Sprintf("Case %d | Mute", len(modlog.Cases)+1),
		Color:     0x00ccff,
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

type CaseUnmute struct {
	ModeratorID string
	UserID      string
	Reason      string
}

func (cm *CaseUnmute) Embed(modlog Modlog, session *discordgo.Session) *discordgo.MessageEmbed {
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
		Title:     fmt.Sprintf("Case %d | Unmute", len(modlog.Cases)+1),
		Color:     0xbb3344,
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
