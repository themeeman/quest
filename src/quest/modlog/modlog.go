package modlog

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"reflect"
	"sync"
	"time"
)

type Modlog struct {
	ChannelID string
	Valid     bool
	Log       chan Case
	Quit      chan struct{}
	Cases     Cases
	Mutex     *sync.Mutex
}

func (m *Modlog) Scan(value interface{}) error {
	*m = Modlog{}
	null := sql.NullString{}
	err := null.Scan(value)
	if err != nil {
		return err
	}
	m.ChannelID = null.String
	m.Valid = null.Valid
	if m.Valid {
		m.Cases = make(Cases, 0)
		m.Log = make(chan Case)
		m.Quit = make(chan struct{})
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
	Embed(*Modlog, *discordgo.Session) *discordgo.MessageEmbed
}

type CaseMessage struct {
	Message string
	Case    Case
}

type Cases []CaseMessage

func (c Cases) MarshalJSON() ([]byte, error) {
	result := make([]interface{}, len(c))
	for i, v := range c {
		t := reflect.TypeOf(v.Case).Elem()
		fields := []reflect.StructField{
			{
				Name: "Message",
				Type: reflect.TypeOf(""),
				Tag:  `json:"message"`,
			},
			{
				Name: "Type",
				Type: reflect.TypeOf(""),
				Tag:  `json:"type"`,
			},
		}
		for i := 0; i < t.NumField(); i++ {
			fields = append(fields, t.Field(i))
		}
		a := reflect.New(reflect.StructOf(fields))
		a.Elem().FieldByName("Message").SetString(v.Message)
		a.Elem().FieldByName("Type").SetString(caseName(v.Case))
		for i := 2; i < a.Elem().NumField(); i++ {
			a.Elem().Field(i).Set(reflect.ValueOf(v.Case).Elem().Field(i - 2))
		}
		result[i] = a.Interface()
	}
	return json.Marshal(result)
}

func (c *Cases) UnmarshalJSON(data []byte) error {
	var temp []map[string]interface{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	*c = make(Cases, len(temp))
	for i, v := range temp {
		(*c)[i].Message = v["message"].(string)
		a := reflect.New(caseType(v["type"].(string)))
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}
		err = json.Unmarshal(d, a.Interface())
		if err != nil {
			return err
		}
		(*c)[i].Case = a.Elem().Interface().(Case)
	}
	return nil
}

func StartLogging(session *discordgo.Session, modlog *Modlog) {
	for {
		select {
		case <-modlog.Quit:
			return
		case c := <-modlog.Log:
			modlog.Mutex.Lock()
			message, err := session.ChannelMessageSendEmbed(modlog.ChannelID, c.Embed(modlog, session))
			if err == nil {
				modlog.Cases = append(modlog.Cases, CaseMessage{
					Message: message.ID,
					Case:    c,
				})
			}
			data, err := modlog.Cases.MarshalJSON()
			fmt.Println(string(data), err)
			modlog.Cases.UnmarshalJSON(data)
			for i, v := range modlog.Cases {
				fmt.Println(i+1, v.Message, v.Case)
			}
			modlog.Mutex.Unlock()
		}
	}
}

func caseName(i interface{}) string {
	switch i.(type) {
	case *CaseBan:
		return "ban"
	case *CaseKick:
		return "kick"
	case *CaseMute:
		return "mute"
	case *CasePurge:
		return "purge"
	case *CaseUnban:
		return "unban"
	case *CaseUnmute:
		return "unmute"
	case *CaseWarn:
		return "warn"
	case *CaseSet:
		return "set"
	}
	return "invalid"
}

func caseType(s string) reflect.Type {
	var a interface{}
	switch s {
	case "ban":
		a = &CaseBan{}
	case "kick":
		a = &CaseKick{}
	case "mute":
		a = &CaseMute{}
	case "purge":
		a = &CasePurge{}
	case "unban":
		a = &CaseUnban{}
	case "unmute":
		a = &CaseUnmute{}
	case "warn":
		a = &CaseWarn{}
	case "set":
		a = &CaseSet{}
	default:
		a = nil
	}
	return reflect.TypeOf(a)
}

func getUser(session *discordgo.Session, id string) (user *discordgo.User) {
	m, _ := session.User(id)
	return m
}

func timeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}
