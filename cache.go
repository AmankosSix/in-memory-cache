package in_memory_cache

import (
	"time"
)

type Cache struct {
	m map[string]interface{}
}

func New() *Cache {
	return &Cache{
		m: make(map[string]interface{}),
	}
}

func (c Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.m[key] = value
	delete := func() {
		c.Delete(key)
	}
	go time.AfterFunc(ttl, delete)
}

func (c Cache) Get(key string) interface{} {
	return c.m[key]
}

func (c Cache) Delete(key string) {
	delete(c.m, key)
}
