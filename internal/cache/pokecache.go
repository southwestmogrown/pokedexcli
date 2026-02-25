package cache

import (
	"sync"
	"time"
)

type entry struct {
    val   []byte
    added time.Time
}

type PokeCache struct {
    mu       sync.RWMutex
    interval time.Duration
    val      map[string]entry
}

func NewPokeCache(interval time.Duration) *PokeCache {
    c := &PokeCache{
        interval: interval,
        val:      make(map[string]entry),
    }
    go c.reapLoop()
    return c
}

func (c *PokeCache) Add(key string, val []byte) {
    cpy := make([]byte, len(val))
    copy(cpy, val)
    c.mu.Lock()
    defer c.mu.Unlock()
    c.val[key] = entry{val: cpy, added: time.Now()}
}

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
