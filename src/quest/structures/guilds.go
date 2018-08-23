package structures

import (
	"../modlog"
	"database/sql"
)

type Guild struct {
	ID            string         `db:"id"`
	MuteRole      sql.NullString `db:"mute_role"      type:"RoleMention"`
	ModRole       sql.NullString `db:"mod_role"       type:"RoleMention"`
	AdminRole     sql.NullString `db:"admin_role"     type:"RoleMention"`
	Modlog        *modlog.Modlog `db:"mod_log"        type:"ChannelMention"`
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

type Guilds map[string]*Guild

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
