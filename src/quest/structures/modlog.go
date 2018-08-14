package structures

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"time"
	"strconv"
	"database/sql/driver"
	"database/sql"
)

type Modlog struct {
	ChannelID string `db:"mod_log"`
	Valid     bool
	Log       chan Case
	Cases     []Case
}

func (m *Modlog) Scan(value interface{}) error {
	null := sql.NullString{}
	err := null.Scan(value)
	fmt.Println(null)
	if err != nil {
		return err
	}
	m.ChannelID = null.String
	m.Valid = null.Valid
	return nil
}

func (m *Modlog) Value() (driver.Value, error) {
	fmt.Println(m)
	null := sql.NullString{
		String: m.ChannelID,
		Valid:  m.Valid,
	}
	return null.Value()
}

type Case interface {
	Embed(*discordgo.Session) *discordgo.MessageEmbed
}

type CaseMute struct {
	ModeratorID string
	UserID      string
	GuildID     string
	Duration    int
	Reason      string
}

func getMember(session *discordgo.Session, guildID string, id string) (member *discordgo.Member) {
	if session.StateEnabled {
		guild, _ := session.State.Guild(guildID)
		if guild != nil {
			member = func() *discordgo.Member {
				for _, m := range guild.Members {
					if m.User.ID == id {
						return m
					}
				}
				return nil
			}()
			if member == nil {
				member, _ = session.GuildMember(guildID, id)
			}
		}
	} else {
		member, _ = session.GuildMember(guildID, id)
	}
	return
}

func timeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

func (cm *CaseMute) Embed(modlog Modlog, session *discordgo.Session) *discordgo.MessageEmbed {
	member := getMember(session, cm.GuildID, cm.UserID)
	moderator := getMember(session, cm.GuildID, cm.UserID)
	return &discordgo.MessageEmbed{
		Title:     fmt.Sprintf("Case %d | Mute", len(modlog.Cases)+1),
		Color:     0x00ccff,
		Timestamp: timeToTimestamp(time.Now().UTC()),
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: moderator.User.AvatarURL(""),
			Name:    moderator.User.String(),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: member.User.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User",
				Value:  member.User.String(),
				Inline: true,
			},
			{
				Name:   "Duration",
				Value:  strconv.Itoa(cm.Duration),
				Inline: true,
			},
			{
				Name:  "Reason",
				Value: cm.Reason,
			},
		},
	}
}
