package protoanalysis

import "sync"

type Cache struct {
	mu   sync.RWMutex
	pkgs map[string]Packages // proto dir path-proto packages pair.
}

func NewCache() *Cache {
	return &Cache{
		pkgs: make(map[string]Packages),
	}
}

func (c *Cache) Get(key string) (Packages, bool) {
	c.mu.RLock()
	pkgs, ok := c.pkgs[key]
	c.mu.RUnlock()
	return pkgs, ok
}

func (c *Cache) Add(key string, value Packages) {
	c.mu.Lock()
	c.pkgs[key] = value
	c.mu.Unlock()
}
