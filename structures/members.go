package structures

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tomvanwoow/quest/inventory"
	"sync"
)

const MemberCacheLimit = 50

type Member struct {
	ID          string                    `db:"user_id"`
	MuteExpires mysql.NullTime            `db:"mute_expires"`
	LastDaily   mysql.NullTime            `db:"last_daily"`
	Experience  int64                     `db:"experience"`
	Chests      inventory.ChestsInventory `db:"chests"`
	*sync.Mutex
}

type Members struct {
	cache Cache
}

func NewMemberCache(db *sqlx.DB, guildID string) Members {
	return Members{
		cache: NewCache(
			db,
			MemberCacheLimit,
			func(db *sqlx.DB, id string) (interface{}, error) {
				return FetchMember(db, guildID, id)
			},
			func(db *sqlx.DB, value interface{}) error {
				return SaveMember(db, guildID, value.(*Member))
			},
			func(id string) interface{} {
				return NewMember(id)
			},
		),
	}
}

func NewMember(id string) *Member {
	return &Member{
		ID:     id,
		Chests: make(inventory.ChestsInventory),
		Mutex:  new(sync.Mutex),
	}
}

func (members Members) Lock() {
	members.cache.mutex.Lock()
}

func (members Members) Unlock() {
	members.cache.mutex.Unlock()
}

func (members Members) Commit(id string) error {
	return members.cache.Commit(id)
}

func (members Members) Add(member *Member) {
	members.cache.Add(member.ID, member)
}

func (members Members) Get(id string) *Member {
	return members.cache.Get(id).(*Member)
}

func (members Members) Destroy(id string) {
	members.cache.Destroy(id)
}

func (members Members) DestroyAll() {
	members.cache.DestroyAll()
}
