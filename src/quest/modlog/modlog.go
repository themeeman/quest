package modlog

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type Modlog struct {
	ChannelID string
	Valid     bool
	Log       chan Case
	Quit      chan struct{}
}

func (m *Modlog) Scan(value interface{}) error {
	if m == nil {
		return nil
	}
	null := sql.NullString{}
	err := null.Scan(value)
	if err != nil {
		return err
	}
	m.ChannelID = null.String
	m.Valid = null.Valid
	if m.Valid {
		m.Log = make(chan Case)
		m.Quit = make(chan struct{})
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

func StartLogging(session *discordgo.Session, modlog Modlog, cases *Cases) {
	for {
		select {
		case <-modlog.Quit:
			return
		case c := <-modlog.Log:
			cases.Mutex.Lock()
			emb := c.Embed(session)
			emb.Title = fmt.Sprintf("Case %d | ", len(cases.Cases)+1) + emb.Title
			message, err := session.ChannelMessageSendEmbed(modlog.ChannelID, emb)
			if err == nil {
				cases.Cases = append(cases.Cases, CaseMessage{
					Message: message.ID,
					Case:    c,
				})
			}
			data, _ := cases.MarshalJSON()
			fmt.Println(string(data), fmt.Sprintf("%p", cases))
			cases.Mutex.Unlock()
		}
	}
}

func getUser(session *discordgo.Session, id string) *discordgo.User {
	m, _ := session.User(id)
	return m
}

func timeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}
