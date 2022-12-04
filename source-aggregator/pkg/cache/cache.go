package cache

import (
	"sync"
	"time"
)

const CacheExpiry = time.Duration(30) * time.Minute

var Cacher = NewCacher()

type cacher struct {
	cacheMap map[string]*CacheEntry
	lock     *sync.RWMutex
}

type CacheEntry struct {
	Value  []byte
	Expiry time.Time
}

func NewCacher() *cacher {
	return &cacher{
		cacheMap: make(map[string]*CacheEntry),
		lock:     &sync.RWMutex{},
	}
}

func (c *cacher) Set(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cacheMap[key] = &CacheEntry{
		Value:  value,
		Expiry: time.Now().Add(CacheExpiry),
	}
}

func (c *cacher) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	entry, ok := c.cacheMap[key]
	if !ok {
		return nil, false
	}

	if entry.Expiry.Before(time.Now()) {
		return nil, false
	}

	return entry.Value, true
}
