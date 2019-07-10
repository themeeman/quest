package modlog

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Modlog struct {
	ChannelID string
	isLogging bool
	log       chan Case
}

func (modlog Modlog) IsLogging() bool {
	return modlog.isLogging
}

func (modlog *Modlog) Scan(value interface{}) error {
	if modlog == nil {
		return nil
	}
	null := sql.NullString{}
	err := null.Scan(value)
	if err != nil {
		return err
	}
	modlog.ChannelID = null.String
	modlog.isLogging = null.Valid
	if modlog.isLogging {
		modlog.log = make(chan Case)
	}
	return nil
}

func (modlog Modlog) Value() (driver.Value, error) {
	null := sql.NullString{
		String: modlog.ChannelID,
		Valid:  modlog.isLogging,
	}
	return null.Value()
}

func (modlog *Modlog) StartLogging(session *discordgo.Session, cases *Cases) {
	modlog.isLogging = true
	for {
		c := <-modlog.log
		cases.Mutex.Lock()
		emb := c.Embed(session)
		emb.Title = fmt.Sprintf("Case %d | ", len(cases.Cases)+1) + emb.Title
		message, err := session.ChannelMessageSendEmbed(modlog.ChannelID, emb)
		if err == nil {
			cases.Cases = append(cases.Cases, &CaseMessage{
				Message: message.ID,
				Case:    c,
			})
		} else {
			fmt.Println(err)
		}
		cases.Mutex.Unlock()
	}
}

func (modlog Modlog) Log(c Case) {
	if modlog.isLogging {
		modlog.log <- c
	}
}