package structures

import (
	"reflect"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"sync"
)

type Case interface {
	Embed(*discordgo.Session, Cases) *discordgo.MessageEmbed
}

type CaseMessage struct {
	Message string
	Case    Case
}

type Cases struct{
	Cases []CaseMessage
	*sync.Mutex
}

func (c Cases) MarshalJSON() ([]byte, error) {
	result := make([]interface{}, len(c.Cases))
	for i, v := range c.Cases {
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
	c.Cases = make([]CaseMessage, len(temp))
	for i, v := range temp {
		c.Cases[i].Message = v["message"].(string)
		a := reflect.New(caseType(v["type"].(string)))
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}
		err = json.Unmarshal(d, a.Interface())
		if err != nil {
			return err
		}
		c.Cases[i].Case = a.Elem().Interface().(Case)
	}
	return nil
}
