package structures

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tomvanwoow/quest/inventory"
	"sync"
	"time"
)

type Member struct {
	ID          string                    `db:"user_id"`
	MuteExpires mysql.NullTime            `db:"mute_expires"`
	LastDaily   mysql.NullTime            `db:"last_daily"`
	Experience  int64                     `db:"experience"`
	Chests      inventory.ChestsInventory `db:"chests"`
	*sync.Mutex
}

type Members struct {
	state map[string]struct {
		lastUsed time.Time
		member   *Member
	}
	db      *sqlx.DB
	mutex   *sync.Mutex
	guildID string
}

func NewMemberCache(db *sqlx.DB, guildID string) Members {
	return Members{
		state: make(map[string]struct {
			lastUsed time.Time
			member   *Member
		}),
		db:      db,
		mutex:   new(sync.Mutex),
		guildID: guildID,
	}
}

func NewMember(id string) *Member {
	return &Member{
		ID:     id,
		Chests: make(inventory.ChestsInventory),
		Mutex:  new(sync.Mutex),
	}
}

func (members Members) getOldest() string {
	var oldestTime time.Time
	var rv string
	for id, g := range members.state {
		if oldestTime.IsZero() || g.lastUsed.Before(oldestTime) {
			oldestTime = g.lastUsed
			rv = id
		}
	}
	return rv
}

func (members *Members) Lock() {
	members.mutex.Lock()
}

func (members *Members) Unlock() {
	members.mutex.Unlock()
}

func (members *Members) Commit(id string) bool {
	members.Lock()
	defer members.Unlock()
	g, ok := members.state[id]
	if !ok {
		return false
	}
	delete(members.state, id)
	err := SaveMember(members.db, members.guildID, g.member)
	if err != nil {
		return false
	}
	return true
}

func (members *Members) CommitAll() {
	members.Lock()
	defer members.Unlock()
	for id := range members.state {
		members.Commit(id)
	}
}

func (members *Members) Add(member *Member) {
	members.Lock()
	defer members.Unlock()
	if len(members.state) >= GuildCacheLimit {
		members.Commit(members.getOldest())
	}
	members.state[member.ID] = struct {
		lastUsed time.Time
		member   *Member
	}{
		lastUsed: time.Now(),
		member:   member,
	}
}

func (members *Members) Get(id string) *Member {
	members.Lock()
	defer members.Unlock()
	if m, ok := members.state[id]; ok {
		m.lastUsed = time.Now()
		return m.member
	}
	member, err := FetchMember(members.db, members.guildID, id)
	if err == nil {
		members.Add(member)
		return member
	}

	member = NewMember(id)
	members.Add(member)
	return member
}

func (members *Members) Destroy(id string) {
	members.Lock()
	defer members.Unlock()
	delete(members.state, id)
}

func (members *Members) DestroyAll() {
	members.Lock()
	defer members.Unlock()
	members.state = make(map[string]struct {
		lastUsed time.Time
		member   *Member
	})
}
