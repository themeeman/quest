package structures

import (
	"../inventory"
	"github.com/go-sql-driver/mysql"
)

type Member struct {
	ID          string                    `db:"user_id"`
	MuteExpires mysql.NullTime            `db:"mute_expires"`
	LastDaily   mysql.NullTime            `db:"last_daily"`
	Experience  int64                     `db:"experience"`
	Chests      inventory.ChestsInventory `db:"chests"`
}

type Members map[string]*Member

func (members Members) Get(id string) *Member {
	member, ok := members[id]
	if !ok {
		member = &Member{ID: id}
		members[id] = member
		return member
	}
	if member.Chests == nil {
		member.Chests = make(inventory.ChestsInventory)
	}
	return member
}
