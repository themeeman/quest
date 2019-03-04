package modlog

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/bwmarrin/discordgo"
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
			fmt.Println(c)
			cases.Mutex.Lock()
			emb := c.Embed(session)
			emb.Title = fmt.Sprintf("Case %d | ", len(cases.Cases)+1) + emb.Title
			message, err := session.ChannelMessageSendEmbed(modlog.ChannelID, emb)
			cm := CaseMessage{
				Message: message.ID,
				Case:    c,
			}
			if err == nil {
				cases.Cases = append(cases.Cases, &cm)
			} else {
				fmt.Println(err)
			}
			data, _ := cm.MarshalJSON()
			fmt.Println(string(data))
			cm = CaseMessage{}
			fmt.Println(cm.UnmarshalJSON(data), cm.Case)
			cases.Mutex.Unlock()
		}
	}
}