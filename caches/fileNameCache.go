package caches

import (
	"sort"
	"speedСontrol/services/config"
	"sync"
)

var Cache *fileNameCache

type fileNameCache struct {
	values    map[string][]int
	hits      map[string]int
	cacheSize int
	mutex     sync.RWMutex
}

func init() {
	Cache = newFileNameCache()
}

func newFileNameCache() *fileNameCache {
	return &fileNameCache{
		values:    make(map[string][]int),
		hits:      make(map[string]int),
		cacheSize: config.Conf.DirsCacheSize,
	}
}

func (c *fileNameCache) Load(key string) (values []int, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	values, ok = c.values[key]

	if val, ok := c.hits[key]; ok {
		c.hits[key] = val + 1
	}

	return
}

func (c *fileNameCache) Store(key string, values []int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.values[key] = values

	if len(c.hits) == c.cacheSize {
		var hits = make([]int, 0, c.cacheSize)

		for _, value := range c.hits {
			if i, ok := contains(hits, value); !ok {
				hits = append(hits, 0)
				copy(hits[i+1:], hits[i:])
				hits[i] = value
			}
		}

		mid := (len(hits) - 1) / 2

		for key, value := range c.hits {
			if value < hits[mid] {
				delete(c.hits, key)
			}
		}
	}

	if _, ok := c.hits[key]; !ok {
		c.hits[key] = 1
	}
}

func contains(s []int, si int) (i int, ok bool) {
	i = sort.SearchInts(s, si)
	ok = i < len(s) && s[i] == si
	return
}
