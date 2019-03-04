package structures

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tomvanwoow/quest/modlog"
	"sync"
	"time"
)

const GuildCacheLimit = 5

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
	Members       Members
	Roles
	*sync.Mutex
}

type fetchFuncGuild func(*sqlx.DB, string) (*Guild, error)
type saveFuncGuild func(*sqlx.DB, *Guild) error

type Guilds struct {
	state map[string]struct {
		lastUsed time.Time
		guild    *Guild
	}
	db    *sqlx.DB
	mutex *sync.Mutex
	fetch fetchFuncGuild
	save  saveFuncGuild
}

func NewGuildCache(db *sqlx.DB, fetch fetchFuncGuild, save saveFuncGuild) Guilds {
	return Guilds{
		state: make(map[string]struct {
			lastUsed time.Time
			guild    *Guild
		}),
		db:    db,
		mutex: new(sync.Mutex),
		fetch: fetch,
		save:  save,
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
		Mutex: new(sync.Mutex),
	}
}

func (guilds Guilds) getOldest() string {
	var oldestTime time.Time
	var rv string
	for id, g := range guilds.state {
		if oldestTime.IsZero() || g.lastUsed.Before(oldestTime) {
			oldestTime = g.lastUsed
			rv = id
		}
	}
	return rv
}

func (guilds *Guilds) Lock() {
	guilds.mutex.Lock()
}

func (guilds *Guilds) Unlock() {
	guilds.mutex.Unlock()
}

func (guilds *Guilds) commitIfLocked(id string) error {
	g, ok := guilds.state[id]
	if !ok {
		return errors.Errorf("guild %s not found", id)
	}
	delete(guilds.state, id)
	err := guilds.save(guilds.db, g.guild)
	if err != nil {
		return err
	}
	return nil
}

func (guilds *Guilds) Commit(id string) error {
	guilds.Lock()
	defer guilds.Unlock()
	return guilds.commitIfLocked(id)
}

func (guilds *Guilds) CommitAll() []error {
	for id, g := range guilds.state {
		g.guild.Members.CommitAll()
		guilds.Commit(id)
	}
	return nil
}

func (guilds *Guilds) addIfLocked(guild *Guild) {
	if len(guilds.state) >= GuildCacheLimit {
		guilds.Commit(guilds.getOldest())
	}
	guilds.state[guild.ID] = struct {
		lastUsed time.Time
		guild    *Guild
	}{
		lastUsed: time.Now(),
		guild:    guild,
	}
}

func (guilds *Guilds) Add(guild *Guild) {
	guilds.Lock()
	defer guilds.Unlock()
	guilds.addIfLocked(guild)
}

func (guilds *Guilds) Get(id string) *Guild {
	guilds.Lock()
	defer guilds.Unlock()
	if g, ok := guilds.state[id]; ok {
		g.lastUsed = time.Now()
		return g.guild
	}
	guild, err := guilds.fetch(guilds.db, id)
	if err == nil {
		guilds.addIfLocked(guild)
		return guild
	}

	guild = NewGuild(id)
	guilds.addIfLocked(guild)
	return guild
}

func (guilds *Guilds) Destroy(id string) {
	guilds.Lock()
	defer guilds.Unlock()
	delete(guilds.state, id)
}

func (guilds *Guilds) DestroyAll() {
	guilds.Lock()
	defer guilds.Unlock()
	guilds.state = make(map[string]struct {
		lastUsed time.Time
		guild    *Guild
	})
}
