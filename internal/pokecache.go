package internal

import (
	"sync"
	"time"
)

type entry struct {
	val   []byte
	added time.Time
}

type PokeCache struct {
	mu        sync.RWMutex // protects the fields below
	interval  time.Duration
	val       map[string]entry
}

func NewPokeCache(interval time.Duration) *PokeCache {
	c := &PokeCache{
		interval: interval,
		val:      make(map[string]entry),
	}
	// start background reaper
	go c.reapLoop()
	return c
}

// Add inserts or replaces a value associated with the given key.
// The cache is protected by a write lock and the stored slice is copied
// so callers cannot mutate the internal representation.
func (c *PokeCache) Add(key string, val []byte) {
	cpy := make([]byte, len(val))
	copy(cpy, val)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val[key] = entry{val: cpy, added: time.Now()}
}

// reapLoop runs in a goroutine and periodically removes entries older
// than the configured interval. The loop exits only when the program
// terminates; there is no shutdown mechanism in this simple cache.
func (c *PokeCache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-c.interval)
		c.mu.Lock()
		for k, e := range c.val {
			if e.added.Before(cutoff) {
				delete(c.val, k)
			}
		}
		c.mu.Unlock()
	}
}

// Get returns a copy of the value for key. The boolean reports whether
// the key was found. It locks for reading.
func (c *PokeCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.val[key]
	if !ok {
		return nil, false
	}
	cpy := make([]byte, len(e.val))
	copy(cpy, e.val)
	return cpy, true
}
