package structures

import (
	"github.com/go-sql-driver/mysql"
	"../inventory"
)

type Member struct {
	ID          string           `db:"user_id"`
	MuteExpires mysql.NullTime   `db:"mute_expires"`
	Experience  int64            `db:"experience"`
	Chests      inventory.Chests `db:"chests"`
}

type Members map[string]*Member

func (members Members) Get(id string) *Member {
	member, ok := members[id]
	if !ok {
		member = &Member{ID: id}
		members[id] = member
		return member
	}
	return member
}
