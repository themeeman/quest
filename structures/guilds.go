package structures

import (
	"../modlog"
	"database/sql"
	"sync"
	"time"
)

type Guild struct {
	ID            string         `db:"id"`
	MuteRole      sql.NullString `db:"mute_role"      type:"RoleMention"`
	ModRole       sql.NullString `db:"mod_role"       type:"RoleMention"`
	AdminRole     sql.NullString `db:"admin_role"     type:"RoleMention"`
	Modlog        modlog.Modlog  `db:"mod_log"        type:"ChannelMention"`
	Autorole      sql.NullString `db:"autorole"       type:"RoleMention"`
	ExpReload     uint16         `db:"exp_reload"     type:"Integer"`
	ExpGainUpper  uint16         `db:"exp_gain_upper" type:"Integer"`
	ExpGainLower  uint16         `db:"exp_gain_lower" type:"Integer"`
	LotteryChance uint8          `db:"lottery_chance" type:"Integer"`
	LotteryUpper  uint32         `db:"lottery_upper"  type:"Integer"`
	LotteryLower  uint32         `db:"lottery_lower"  type:"Integer"`
	Cases         modlog.Cases   `db:"cases"`
	Members
	Roles
}

type Guilds struct {
	state    map[string]*Guild
	lastUsed map[string]time.Time
	*sync.Mutex
}

func NewGuildCache() Guilds {
	return Guilds{
		state:    make(map[string]*Guild),
		lastUsed: make(map[string]time.Time),
		Mutex:    new(sync.Mutex),
	}
}

func NewGuild(id string) *Guild {
	return &Guild{
		ID:            id,
		ExpReload:     60,
		ExpGainUpper:  25,
		ExpGainLower:  10,
		LotteryChance: 100,
		LotteryUpper:  500,
		LotteryLower:  250,
		Cases: modlog.Cases{
			Cases: make([]*modlog.CaseMessage, 0, 1000),
		},
		Modlog: modlog.Modlog{
			Log: make(chan modlog.Case),
		},
	}
}

func (guilds Guilds) getOldest() string {
	var oldestTime time.Time
	var rv string
	for id, t := range guilds.lastUsed {
		if oldestTime.IsZero() || t.Before(oldestTime) {
			oldestTime = t
			rv = id
		}
	}
	return rv
}

func (guilds Guilds) removeOldest() {
	delete(guilds.state, guilds.getOldest())
}

func (guilds *Guilds) Get(id string) *Guild {
	if guilds == nil {
		return nil
	}
	guild, ok := guilds.state[id]
	if !ok {
		guilds.state[id] = NewGuild(id)
	}
	if guild.Members == nil {
		guild.Members = make(Members)
	}
	return guild
}
