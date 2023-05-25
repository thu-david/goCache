package caches

import (
	"goCache/helpers"
	"sync"
)

type Cache struct{

	data map[string][]byte
	cnt int64
	lock *sync.RWMutex
}

func NewCache() *Cache{
	return &Cache{
		data: make(map[string][]byte, 256),
		cnt: 0,
		lock: &sync.RWMutex{},
	}
}

func (c* Cache)Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c* Cache)Set(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_ , ok := c.data[key]
	if !ok {
		c.cnt++
	}
	c.data[key] = helpers.Copy(value)
}

func (c* Cache)Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	if ok {
		c.cnt--
		delete(c.data, key)
	}
}

func (c* Cache)Count() int64{
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cnt
}