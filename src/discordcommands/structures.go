package discordcommands

import (
	"github.com/bwmarrin/discordgo"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"database/sql/driver"
	"time"
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
	context Bot) BotError

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
	ID        string         `db:"id"`
	MuteRole  sql.NullString `db:"mute_role"  type:"RoleMention"`
	ModRole   sql.NullString `db:"mod_role"   type:"RoleMention"`
	AdminRole sql.NullString `db:"admin_role" type:"RoleMention"`
	Modlog    sql.NullString `db:"mod_log"    type:"RoleMention"`
	Autorole  sql.NullString `db:"autorole"   type:"RoleMention"`
	ExpReload uint16         `db:"exp_reload" type:"Integer"`
	Members
	Roles
}

type Member struct {
	ID          string   `db:"user_id"`
	MuteExpires NullTime `db:"mute_expires"`
	Experience  int64    `db:"experience"`
}

type Role struct {
	Experience int64  `db:"experience"`
	ID         string `db:"id"`
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
		guild = &Guild{ID: id}
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

//All NullTime code stolen from	lib/pq
type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (d Command) ForcedArgs() (i int) {
	for _, v := range d.Arguments {
		if !v.Optional {
			i += 1
		}
	}
	return
}

func MustGetGuildID(session *discordgo.Session, message *discordgo.MessageCreate) string {
	must := func(v *discordgo.Channel, _ error) *discordgo.Channel { return v }
	return must(session.Channel(message.ChannelID)).GuildID
}
