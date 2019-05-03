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
		lastUsed int64
		value    interface{}
	}
	db                 *sqlx.DB
	mutex              *sync.Mutex
	capacity           int
	fetch              fetchFunc
	save               saveFunc
	newValue           newValueFunc
	commitErrorMessage func(id string) string
}

func NewCache(db *sqlx.DB, capacity int, fetch fetchFunc, save saveFunc, newValue newValueFunc, commitError func(id string) string) Cache {
	return Cache{
		state: make(map[string]struct {
			lastUsed int64
			value    interface{}
		}, capacity),
		db:                 db,
		mutex:              new(sync.Mutex),
		capacity:           capacity,
		fetch:              fetch,
		save:               save,
		newValue:           newValue,
		commitErrorMessage: commitError,
	}
}

func (cache Cache) getOldest() string {
	var oldestTime int64
	var rv string
	for id, g := range cache.state {
		if oldestTime == 0 || g.lastUsed < oldestTime {
			oldestTime = g.lastUsed
			rv = id
		}
	}
	return rv
}

func (cache *Cache) lock() {
	cache.mutex.Lock()
}

func (cache *Cache) unlock() {
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
	cache.lock()
	defer cache.unlock()
	return cache.commitIfLocked(id)
}

func (cache *Cache) CommitAll() []error {
	cache.lock()
	defer cache.unlock()
	errs := make([]error, 0, cache.capacity)
	for id := range cache.state {
		err := cache.commitIfLocked(id)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "error commiting member %s:", id))
		}
	}
	return errs
}

func (cache *Cache) addIfLocked(id string, value interface{}) {
	if len(cache.state) >= cache.capacity {
		_ = cache.Commit(cache.getOldest())
	}
	cache.state[id] = struct {
		lastUsed int64
		value    interface{}
	}{
		lastUsed: time.Now().Unix(),
		value:    value,
	}
}

func (cache *Cache) Add(id string, value interface{}) {
	cache.lock()
	defer cache.unlock()
	cache.addIfLocked(id, value)
}

func (cache *Cache) Get(id string) interface{} {
	cache.lock()
	defer cache.unlock()
	if g, ok := cache.state[id]; ok {
		g.lastUsed = time.Now().Unix()
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
	cache.lock()
	defer cache.unlock()
	delete(cache.state, id)
}

func (cache *Cache) DestroyAll() {
	cache.lock()
	defer cache.unlock()
	cache.state = make(map[string]struct {
		lastUsed int64
		value    interface{}
	})
}
