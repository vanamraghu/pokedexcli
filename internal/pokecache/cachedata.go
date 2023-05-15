package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Mux       *sync.Mutex
	CacheData map[string]CacheEntry
}

// Add Adding the key
func (c *Cache) Add(key string, value []byte) {
	c.Mux.Lock()
	c.CacheData[key] = CacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	defer c.Mux.Unlock()
}

// Get Getting the value
func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	value, ok := c.CacheData[key]
	//if ok {
	//	fmt.Printf("%s\n", value.val)
	//} else {
	//	fmt.Println("Has to request from api")
	//}
	return value.val, ok
}

// NewCache Initialize the CacheData by calling this NewCache
func NewCache(interval time.Duration) Cache {
	c := Cache{
		Mux:       &sync.Mutex{},
		CacheData: make(map[string]CacheEntry),
	}
	go c.ReapLoop(interval)
	return c
}

// ReapLoop Reap loop which would be called by NewCache
func (c *Cache) ReapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(now time.Time, last time.Duration) {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	for k, v := range c.CacheData {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.CacheData, k)
		}
	}
}
