package cache

// MemCache is in-memory cache backed by map[string]string
type MemCache struct {
	items map[string]string
}

// Get a value for key k
func (c *MemCache) Get(k string) (string, bool) {
	v, exist := c.items[k]
	return v, exist
}

// Set value v for key k
func (c *MemCache) Set(k string, v string) {
	c.items[k] = v
}

// Returns a new cache
func New() *MemCache {
	m := make(map[string]string)
	return &MemCache{
		items: m,
	}
}
