package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"reflect"
	"bytes"
	"fmt"
	"strings"
)

type Argument struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Infinite bool   `json:"infinite"`
}

type Command struct {
	Description string      `json:"description"`
	Arguments   []*Argument `json:"arguments"`
	Cooldown    int         `json:"cooldown"`
	Group                   `json:"permission"`
	Aliases     []string    `json:"aliases"`
	Examples    []string    `json:"examples"`
	Hidden      bool        `json:"hidden"`
}

type Handler func(session *discordgo.Session,
	message *discordgo.MessageCreate,
	args map[string]string) error

type HandlerMap map[string]Handler

type CommandMap map[string]*Command

func (c Command) ForcedArgs() (i int) {
	for _, v := range c.Arguments {
		if !v.Optional {
			i += 1
		}
	}
	return
}

func (c Command) GetUsage(prefix string, name string) string {
	var buffer bytes.Buffer
	buffer.WriteString(prefix + name)
	for _, v := range c.Arguments {
		if v.Optional {
			buffer.WriteString(fmt.Sprintf(" <%s>", v.Name))
		} else {
			buffer.WriteString(fmt.Sprintf(" [%s]", v.Name))
		}
	}
	return buffer.String()
}

func MustGetGuildID(session *discordgo.Session, message *discordgo.MessageCreate) string {
	c, _ := session.Channel(message.ChannelID)
	if c != nil {
		return c.GuildID
	} else {
		return MustGetGuildID(session, message)
	}
}

func Contains(slice interface{}, value interface{}) (bool, int) {
	s := reflect.ValueOf(slice)
	if !(s.Kind() == reflect.Slice || s.Kind() == reflect.Array) {
		panic("Slice must be a slice!")
	}
	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(value, s.Index(i).Interface()) {
			return true, i
		}
	}
	return false, 0
}

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.ToLower(s)[0:len(prefix)] == strings.ToLower(prefix)
}

func TrimPrefix(s, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}
