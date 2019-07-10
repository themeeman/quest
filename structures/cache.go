package structures

import (
	"sync"
	"time"
)

// Rewrite this with generics when Go2 is released

type fetchFunc func(string) (sync.Locker, error)
type saveFunc func(sync.Locker) error
type newValueFunc func(string) sync.Locker

type Cache struct {
	state map[string]struct {
		lastUsed int64
		value    sync.Locker
	}
	mutex       *sync.Mutex
	capacity    int
	fetch       fetchFunc
	save        saveFunc
	newValue    newValueFunc
	commitError func(id string) error
}

func NewCache(capacity int, fetch fetchFunc, save saveFunc, newValue newValueFunc, commitError func(id string) error) Cache {
	return Cache{
		state: make(map[string]struct {
			lastUsed int64
			value    sync.Locker
		}, capacity),
		mutex:       new(sync.Mutex),
		capacity:    capacity,
		fetch:       fetch,
		save:        save,
		newValue:    newValue,
		commitError: commitError,
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
		return cache.commitError(id)
	}
	delete(cache.state, id)
	err := cache.save(g.value)
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
			errs = append(errs, cache.commitError(id))
		}
	}
	return errs
}

func (cache *Cache) addIfLocked(id string, value sync.Locker) {
	if len(cache.state) >= cache.capacity {
		_ = cache.Commit(cache.getOldest())
	}
	cache.state[id] = struct {
		lastUsed int64
		value    sync.Locker
	}{
		lastUsed: time.Now().Unix(),
		value:    value,
	}
}

func (cache *Cache) Add(id string, value sync.Locker) {
	cache.lock()
	defer cache.unlock()
	cache.addIfLocked(id, value)
}

func (cache *Cache) Get(id string) sync.Locker {
	cache.lock()
	defer cache.unlock()
	if g, ok := cache.state[id]; ok {
		g.lastUsed = time.Now().Unix()
		return g.value
	}
	value, err := cache.fetch(id)
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
		value    sync.Locker
	})
}

func (cache *Cache) Apply(id string, f func(sync.Locker)) {
	value := cache.Get(id)
	value.Lock()
	defer value.Unlock()
	f(value)
}
