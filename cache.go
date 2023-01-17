package in_memory_cache

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

//var timer *time.Timer

type Cache struct {
	m   map[string]interface{}
	mu  *sync.Mutex
	ttl *time.Timer
}

func New() *Cache {
	return &Cache{
		m:   make(map[string]interface{}),
		mu:  new(sync.Mutex),
		ttl: nil,
	}
}

func (c Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	c.m[key] = value
	if c.ttl == nil {
		c.ttl = time.NewTimer(ttl)
		go c.startTimer(key)
	} else {
		c.ttl.Stop()
		c.ttl.Reset(ttl)
	}
	c.mu.Unlock()
}

func (c Cache) startTimer(key string) {
	for {
		select {
		case <-c.ttl.C:
			log.Println("Time to clear cache")
			c.Delete(key)
		default:
			log.Println("Waiting")
			time.Sleep(1 * time.Second)
		}
	}
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
