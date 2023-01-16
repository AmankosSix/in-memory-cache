package in_memory_cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var timer *time.Timer

type Cache struct {
	m  map[string]interface{}
	mu *sync.Mutex
}

func New() *Cache {
	return &Cache{
		m:  make(map[string]interface{}),
		mu: new(sync.Mutex),
	}
}

func (c Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	c.m[key] = value
	timer.Reset(ttl)
	c.mu.Unlock()

	go func() {
		<-timer.C
		c.Delete(key)
	}()
}

func (c Cache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.m[key]
	if ok {
		return val, nil
	}
	message := fmt.Sprintf("%s is not Found in cache memory", key)
	return nil, errors.New(message)
}

func (c Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.m, key)
	c.mu.Unlock()
}
