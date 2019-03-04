package modlog

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"reflect"
	"sync"
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
)

type Case interface {
	Embed(*discordgo.Session) *discordgo.MessageEmbed
}

type CaseMessage struct {
	Message string
	Case    Case
}

type Cases struct {
	Cases []*CaseMessage
	sync.Mutex
}

func (c *CaseMessage) MarshalJSON() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	t := reflect.TypeOf(c.Case).Elem()
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
	a.Elem().FieldByName("Message").SetString(c.Message)
	a.Elem().FieldByName("Type").SetString(caseName(c.Case))
	for i := 2; i < a.Elem().NumField(); i++ {
		a.Elem().Field(i).Set(reflect.ValueOf(c.Case).Elem().Field(i - 2))
	}
	d, err := json.Marshal(a.Interface())
	return d, errors.WithStack(err)
}

func (c *CaseMessage) UnmarshalJSON(data []byte) error {
	var temp map[string]interface{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return errors.WithStack(err)
	}
	c.Message = temp["message"].(string)
	T := caseType(temp["type"].(string))
	if T == nil {
		return nil
	}
	strucct := reflect.New(T)
	d, err := json.Marshal(temp)
	if err != nil {
		return errors.WithStack(err)
	}
	err = json.Unmarshal(d, strucct.Interface())
	if err != nil {
		return errors.WithStack(err)
	}
	c.Case = strucct.Elem().Interface().(Case)
	return nil
}

func (c Cases) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Cases)
}

func (c *Cases) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &c.Cases)
	return err
}

func (c *Cases) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		err := json.Unmarshal(v, c)
		if err != nil {
			return err
		}
		return nil
	case string:
		err := json.Unmarshal([]byte(v), c)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (c Cases) Value() (driver.Value, error) {
	return json.Marshal(c)
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
	case *CaseAddExp:
		return "addexp"
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
	case "addexp":
		a = &CaseAddExp{}
	default:
		a = nil
	}
	return reflect.TypeOf(a)
}
