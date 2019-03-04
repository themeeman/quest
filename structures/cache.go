package structures

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type fetchFunc func(*sqlx.DB, string) (interface{}, error)
type saveFunc func(*sqlx.DB, interface{}) error

type Cache struct {
	state map[string]struct {
		lastUsed time.Time
		value    interface{}
	}
	db    *sqlx.DB
	mutex *sync.Mutex
	fetch fetchFunc
	save  saveFunc
}

func NewCache(db *sqlx.DB, fetch fetchFunc, save saveFunc) Cache {
	return Cache{
		state: make(map[string]struct {
			lastUsed time.Time
			value    interface{}
		}),
		db:    db,
		mutex: new(sync.Mutex),
		fetch: fetch,
		save:  save,
	}
}

func (cache Cache) getOldest() string {
	var oldestTime time.Time
	var rv string
	for id, g := range cache.state {
		if oldestTime.IsZero() || g.lastUsed.Before(oldestTime) {
			oldestTime = g.lastUsed
			rv = id
		}
	}
	return rv
}

func (cache *Cache) Lock() {
	cache.mutex.Lock()
}

func (cache *Cache) Unlock() {
	cache.mutex.Unlock()
}

func (cache *Cache) commitIfLocked(id string) error {
	g, ok := cache.state[id]
	if !ok {
		return errors.Errorf("value %s not found", id)
	}
	delete(cache.state, id)
	err := cache.save(cache.db, g.value)
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) Commit(id string) error {
	cache.Lock()
	defer cache.Unlock()
	return cache.commitIfLocked(id)
}

func (cache *Cache) CommitAll() []error {
	for id := range cache.state {
		cache.Commit(id)
	}
	return nil
}

func (cache *Cache) addIfLocked(id string, value interface{}) {
	if len(cache.state) >= GuildCacheLimit {
		cache.Commit(cache.getOldest())
	}
	cache.state[id] = struct {
		lastUsed time.Time
		value    interface{}
	}{
		lastUsed: time.Now(),
		value:    value,
	}
}

func (cache *Cache) Add(id string, value interface{}) {
	cache.Lock()
	defer cache.Unlock()
	cache.addIfLocked(id, value)
}

func (cache *Cache) Get(id string) interface{} {
	cache.Lock()
	defer cache.Unlock()
	if g, ok := cache.state[id]; ok {
		g.lastUsed = time.Now()
		return g.value
	}
	value, err := cache.fetch(cache.db, id)
	if err == nil {
		cache.addIfLocked(id, value)
		return value
	}

	value = NewGuild(id)
	cache.addIfLocked(id, value)
	return value
}

func (cache *Cache) Destroy(id string) {
	cache.Lock()
	defer cache.Unlock()
	delete(cache.state, id)
}

func (cache *Cache) DestroyAll() {
	cache.Lock()
	defer cache.Unlock()
	cache.state = make(map[string]struct {
		lastUsed time.Time
		value    interface{}
	})
}
