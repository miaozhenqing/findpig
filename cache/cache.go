package cache

import (
	"findpig/util"
	"time"
)

type Item struct {
	Value    any
	ExpireAt int64
}

type Cache struct {
	items         map[string]Item
	cleanInterval time.Duration
	lastCleanTime time.Time
}

func NewCache(cleanInterval time.Duration) *Cache {
	c := &Cache{
		items:         make(map[string]Item),
		cleanInterval: cleanInterval,
		lastCleanTime: util.CurrentTime(),
	}
	return c
}
func (c *Cache) Get(key string) (any, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if item.ExpireAt < 0 || item.ExpireAt > util.CurrentTime().UnixMilli() {
		return item.Value, true
	}
	return nil, false
}
func (c *Cache) Set(key string, value any, expire int64) {
	go c.startCleanExpired()
	c.items[key] = Item{
		Value:    value,
		ExpireAt: util.CurrentTime().UnixMilli() + expire,
	}
}
func (c *Cache) startCleanExpired() {
	if util.CurrentTime().Sub(c.lastCleanTime) < c.cleanInterval {
		return
	}
	cnt := 0
	for key, value := range c.items {
		if value.ExpireAt < util.CurrentTime().UnixMilli() {
			delete(c.items, key)
			cnt++
		}
		if cnt > 10 {
			break
		}
	}
}
