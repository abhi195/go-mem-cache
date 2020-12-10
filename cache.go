package cache

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultItemTTL         time.Duration = 30 * time.Minute
	defaultCleanupInterval time.Duration = 60 * time.Minute
)

var (
	logger *log.Logger
)

func init() {
	logger = log.New()
	logger.SetFormatter(&log.JSONFormatter{})
}

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
		logger.Tracef("key=%s not found", k)
		mc.mu.RUnlock()
		return "", false
	}
	if time.Now().UnixNano() > i.ttl {
		logger.Tracef("key=%s is expired", k)
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
	logger.Trace("Expired key eviction in progress...")
	totalEvicted := 0
	for k, v := range mc.items {
		if now > v.ttl {
			logger.Tracef("Evicting key=%s", k)
			mc.delete(k)
			totalEvicted++
		}
	}
	logger.Tracef("Evicted %d expired keys", totalEvicted)
	mc.mu.Unlock()
}

type cleaner struct {
	interval time.Duration
}

func (c *cleaner) Run(mc *MemCache) {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		mc.evictExpired()
	}
}

func runCleaner(mc *MemCache, ci time.Duration) {
	c := &cleaner{
		interval: ci,
	}

	mc.cleaner = c
	go c.Run(mc)
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
	logger.Infof("Creating cache with ttl=%s and cleanupInterval=%s", itemTTL, cleanupInterval)
	mc := newCache(itemTTL)
	runCleaner(mc, cleanupInterval)
	return mc
}
