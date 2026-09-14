package redis

import (
	"sync"
	"time"
)

// DedupCache tracks seen message IDs to prevent duplicate processing
type DedupCache struct {
	seen sync.Map
}

var GlobalDedupCache = &DedupCache{}

func init() {
	go GlobalDedupCache.pruneLoop()
}

// CheckAndAdd returns true if the id was already seen, false otherwise
// If false, it adds the id to the cache.
func (c *DedupCache) CheckAndAdd(id string) bool {
	if id == "" {
		return false
	}
	
	_, loaded := c.seen.LoadOrStore(id, time.Now())
	return loaded
}

func (c *DedupCache) pruneLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		c.seen.Range(func(key, value interface{}) bool {
			timestamp := value.(time.Time)
			if now.Sub(timestamp) > 5*time.Minute {
				c.seen.Delete(key)
			}
			return true
		})
	}
}
