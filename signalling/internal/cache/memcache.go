package cache

var MemcacheInstance = NewMemcache()

type Memcache struct {
	cache map[string]interface{}
}

func NewMemcache() *Memcache {
	return &Memcache{
		cache: make(map[string]interface{}),
	}
}

func (c *Memcache) Set(key string, value interface{}) {
	c.cache[key] = value
}

func (c *Memcache) Get(key string) (interface{}, bool) {
	value, ok := c.cache[key]
	return value, ok
}

func (c *Memcache) Delete(key string) {
	delete(c.cache, key)
}

func (c *Memcache) Contains(key string) bool {
	_, ok := c.cache[key]
	return ok
}
