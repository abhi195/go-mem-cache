package cache

import (
	"runtime"
	"sync"
	"time"
)

const (
	defaultItemTTL         time.Duration = 30 * time.Minute
	defaultCleanupInterval time.Duration = 60 * time.Minute
)

type item struct {
	value string
	ttl   int64
}

// MemCache is in-memory cache
type MemCache struct {
	itemTTL time.Duration
	items   map[string]item
	mu      sync.RWMutex
	cleaner *cleaner
}

// Get a value for key k
func (mc *MemCache) Get(k string) (string, bool) {
	mc.mu.RLock()
	i, exist := mc.items[k]
	if !exist {
		mc.mu.RUnlock()
		return "", false
	}
	if time.Now().UnixNano() > i.ttl {
		mc.mu.RUnlock()
		return "", false
	}
	mc.mu.RUnlock()
	return i.value, true
}

// Set value v for key k
func (mc *MemCache) Set(k string, v string) {
	e := time.Now().Add(mc.itemTTL).UnixNano()
	mc.mu.Lock()
	mc.items[k] = item{
		value: v,
		ttl:   e,
	}
	mc.mu.Unlock()
}

func (mc *MemCache) delete(k string) {
	delete(mc.items, k)
}

func (mc *MemCache) evictExpired() {
	now := time.Now().UnixNano()
	mc.mu.Lock()
	for k, v := range mc.items {
		if now > v.ttl {
			mc.delete(k)
		}
	}
	mc.mu.Unlock()
}

type cleaner struct {
	interval time.Duration
	stop     chan bool
}

func (c *cleaner) Run(mc *MemCache) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			mc.evictExpired()
		}
	}
}

func runCleaner(mc *MemCache, ci time.Duration) {
	c := &cleaner{
		interval: ci,
		stop:     make(chan bool),
	}

	mc.cleaner = c
	go c.Run(mc)
}

func stopCleaner(mc *MemCache) {
	mc.cleaner.stop <- true
}

func newCache(ttl time.Duration) *MemCache {
	m := make(map[string]item)
	return &MemCache{
		itemTTL: ttl,
		items:   m,
	}
}

// Returns new MemCache with default itemTTL and cleanupInterval
func New() *MemCache {
	return NewWith(defaultItemTTL, defaultCleanupInterval)
}

// Returns new MemCache with mentioned itemTTL and cleanupInterval
func NewWith(itemTTL time.Duration, cleanupInterval time.Duration) *MemCache {
	mc := newCache(itemTTL)
	runCleaner(mc, cleanupInterval)
	runtime.SetFinalizer(mc, stopCleaner)
	return mc
}
