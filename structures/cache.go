package structures

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"sync"
	"time"
)

// Rewrite this with generics when Go2 is released

type fetchFunc func(*sqlx.DB, string) (interface{}, error)
type saveFunc func(*sqlx.DB, interface{}) error
type newValueFunc func(string) interface{}

type Cache struct {
	state map[string]struct {
		lastUsed time.Time
		value    interface{}
	}
	db       *sqlx.DB
	mutex    *sync.Mutex
	limit    int
	fetch    fetchFunc
	save     saveFunc
	newValue newValueFunc
}

func NewCache(db *sqlx.DB, limit int, fetch fetchFunc, save saveFunc, newValue newValueFunc) Cache {
	return Cache{
		state: make(map[string]struct {
			lastUsed time.Time
			value    interface{}
		}),
		db:       db,
		mutex:    new(sync.Mutex),
		limit:    limit,
		fetch:    fetch,
		save:     save,
		newValue: newValue,
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
	cache.Lock()
	defer cache.Unlock()
	errs := make([]error, 0, cache.limit)
	for id := range cache.state {
		errs = append(errs, cache.commitIfLocked(id))
	}
	return errs
}

func (cache *Cache) addIfLocked(id string, value interface{}) {
	if len(cache.state) >= cache.limit {
		_=cache.Commit(cache.getOldest())
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

	value = cache.newValue(id)
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
