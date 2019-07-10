package structures

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tomvanwoow/quest/modlog"
	"sync"
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
	*sync.RWMutex
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
		RWMutex: new(sync.RWMutex),
	}
}

type Guilds struct {
	cache Cache
}

func NewGuildsCache(db *sqlx.DB) Guilds {
	return Guilds{
		cache: NewCache(
			GuildCacheLimit,
			func(id string) (sync.Locker, error) {
				return FetchGuild(db, id)
			},
			func(value sync.Locker) error {
				return SaveGuild(db, value.(*Guild))
			},
			func(id string) sync.Locker {
				return NewGuild(id)
			},
			func(id string) error {
				return errors.Errorf("error committing guild %s: ", id)
			},
		),
	}
}

func (guilds Guilds) Commit(id string) error {
	return guilds.cache.Commit(id)
}

func (guilds Guilds) CommitAll() []error {
	return guilds.cache.CommitAll()
}

func (guilds Guilds) Add(guild *Guild) {
	guilds.cache.Add(guild.ID, guild)
}

func (guilds Guilds) Get(id string) *Guild {
	return guilds.cache.Get(id).(*Guild)
}

func (guilds Guilds) Destroy(id string) {
	guilds.cache.Destroy(id)
}

func (guilds Guilds) DestroyAll() {
	guilds.cache.DestroyAll()
}

func (guilds Guilds) Apply(id string, f func(*Guild)) {
	guilds.cache.Apply(id, func(value sync.Locker) {
		f(value.(*Guild))
	})
}

func CommitAllGuilds(guilds Guilds) []error {
	errs := make([]error, 0, GuildCacheLimit * (MemberCacheLimit + 1))
	for id := range guilds.cache.state {
		errs = append(errs, guilds.Get(id).Members.CommitAll()...)
	}
	errs = append(errs, guilds.CommitAll()...)
	return errs
}