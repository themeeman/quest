package discordcommands

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"reflect"
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
	Group       `json:"permission"`
	Aliases     []string `json:"aliases"`
	Examples    []string `json:"examples"`
	Hidden      bool     `json:"hidden"`
}

type Handler func(session *discordgo.Session,
	message *discordgo.MessageCreate,
	args map[string]string) error

type HandlerMap map[string]Handler

type CommandMap map[string]*Command

type Group uint

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
func (c Command) ForcedArgs() (i int) {
	for _, v := range c.Arguments {
		if !v.Optional {
			i += 1
		}
	}
	return
}

func (c Command) GetUsage(prefix, name string) string {
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

func Contains(slice, value interface{}) (bool, int) {
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
