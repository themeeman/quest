package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
	"reflect"
	"github.com/go-sql-driver/mysql"
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
	Permission  int         `json:"permission"`
	Aliases     []string    `json:"aliases"`
	Examples    []string    `json:"examples"`
	Hidden      bool        `json:"hidden"`
}

type Handler func(session *discordgo.Session,
	message *discordgo.MessageCreate,
	args map[string]string,
	bot *Bot) BotError

type HandlerMap map[string]Handler

type CommandMap map[string]*Command

type Bot struct {
	HandlerMap
	CommandMap
	ExpTimes map[struct {
		Guild  string
		Member string
	}]time.Time
	Regex  map[string]string
	Prefix string
	Guilds
	DB     *sqlx.DB
	Embed func(title string,
		description string,
		fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed
}

type Guild struct {
	ID            string         `db:"id"`
	MuteRole      sql.NullString `db:"mute_role"      type:"RoleMention"`
	ModRole       sql.NullString `db:"mod_role"       type:"RoleMention"`
	AdminRole     sql.NullString `db:"admin_role"     type:"RoleMention"`
	Modlog        sql.NullString `db:"mod_log"        type:"RoleMention"`
	Autorole      sql.NullString `db:"autorole"       type:"RoleMention"`
	ExpReload     uint16         `db:"exp_reload"     type:"Integer"`
	ExpGainUpper  uint16         `db:"exp_gain_upper" type:"Integer"`
	ExpGainLower  uint16         `db:"exp_gain_lower" type:"Integer"`
	LotteryChance uint8          `db:"lottery_chance" type:"Integer"`
	LotteryUpper  uint32         `db:"lottery_upper"  type:"Integer"`
	LotteryLower  uint32         `db:"lottery_lower"  type:"Integer"`
	Members
	Roles
}

type Member struct {
	ID          string         `db:"user_id"`
	MuteExpires mysql.NullTime `db:"mute_expires"`
	Experience  int64          `db:"experience"`
}

type Role struct {
	ID         string `db:"id"`
	Experience int64  `db:"experience"`
}

type Option struct {
	Name string
	Type string
}

type Guilds map[string]*Guild
type Members map[string]*Member

type Roles []*Role

func (r Roles) Len() int           { return len(r) }
func (r Roles) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Roles) Less(i, j int) bool { return r[i].Experience < r[j].Experience }

func (guilds Guilds) Get(id string) *Guild {
	guild, ok := guilds[id]
	if !ok {
		guild = &Guild{
			ID:            id,
			ExpReload:     60,
			ExpGainUpper:  25,
			ExpGainLower:  10,
			LotteryChance: 100,
			LotteryUpper:  500,
			LotteryLower:  250,
		}
		guilds[id] = guild
	}
	if guild.Members == nil {
		guild.Members = make(Members)
	}
	return guild
}

func (members Members) Get(id string) *Member {
	member, ok := members[id]
	if !ok {
		member = &Member{ID: id}
		members[id] = member
		return member
	}
	return member
}

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

func IsDirectMessage(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	c, _ := session.Channel(message.ChannelID)
	if c != nil {
		g, _ := session.Guild(c.GuildID)
		return g == nil
	} else {
		return IsDirectMessage(session, message)
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
