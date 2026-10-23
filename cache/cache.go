// Package cache holds a small LRU cache.
package cache

import (
	"container/list"
	"sync"
)

type entry struct{ key, val string }

// Cache is an LRU cache that is safe for concurrent use.
type Cache struct {
	mu    sync.Mutex
	cap   int
	order *list.List
	items map[string]*list.Element
}

// New returns a cache that holds at most capacity entries.
func New(capacity int) *Cache {
	return &Cache{cap: capacity, order: list.New(), items: map[string]*list.Element{}}
}

// Get returns the value of key and whether it was present.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.items[key]
	if !ok {
		return "", false
	}
	return el.Value.(*entry).val, true
}

// Put stores val under key and evicts the least recently used entry when the
// cache is full.
func (c *Cache) Put(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		el.Value.(*entry).val = val
		c.order.MoveToFront(el)
		return
	}
	c.items[key] = c.order.PushFront(&entry{key, val})
	if c.order.Len() > c.cap {
		last := c.order.Back()
		c.order.Remove(last)
		delete(c.items, last.Value.(*entry).key)
	}
}

// Len returns the number of entries.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

// GetOrLoad returns the value of key, calling load to make it when absent.
func (c *Cache) GetOrLoad(key string, load func() (string, error)) (string, error) {
	if v, ok := c.Get(key); ok {
		return v, nil
	}
	v, err := load()
	if err != nil {
		return "", err
	}
	c.Put(key, v)
	return v, nil
}
