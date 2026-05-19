package main

import (
	"sync"
	"time"
)

// ttlCache is a generic-ish TTL cache for slow-changing data such as the
// stats table, alerts table, and astro calculations.
//
// Each call to /api/current or every WebSocket tick previously invoked
// getConditions() which calls minmax() ~10 times, each running a full
// SELECT * FROM stats. With a 60s TTL we collapse that to one SELECT per
// minute regardless of request volume.
type ttlCache struct {
	mu    sync.RWMutex
	value interface{}
	at    time.Time
	ttl   time.Duration
}

func newTTLCache(ttl time.Duration) *ttlCache {
	return &ttlCache{ttl: ttl}
}

// Get returns the cached value if it's still fresh. Otherwise it calls
// fn() to refresh, stores, and returns the new value. fn is only called
// inside the write lock so concurrent misses don't stampede the DB.
func (c *ttlCache) Get(fn func() interface{}) interface{} {
	c.mu.RLock()
	if c.value != nil && time.Since(c.at) < c.ttl {
		v := c.value
		c.mu.RUnlock()
		return v
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	// Re-check under the write lock in case another goroutine refreshed
	// while we were waiting.
	if c.value != nil && time.Since(c.at) < c.ttl {
		return c.value
	}
	c.value = fn()
	c.at = time.Now()
	return c.value
}

// Invalidate clears the cache (useful for tests or manual refresh).
func (c *ttlCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = nil
}

// Package-level caches for the slow-changing data.
var (
	statsCache  = newTTLCache(60 * time.Second)
	alertsCache = newTTLCache(60 * time.Second)
)
