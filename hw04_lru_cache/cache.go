package hw04lrucache

import "sync"

type Key string

type cachedItem struct {
	k Key
	v interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    *sync.Mutex
}

func (c *lruCache) GCLast() {
	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back()
		delete(c.items, lastItem.Value.(cachedItem).k)
		c.queue.Remove(lastItem)
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, isSet := c.items[key]
	if isSet {
		c.queue.Remove(item)
	}
	c.items[key] = c.queue.PushFront(cachedItem{k: key, v: value})

	c.GCLast()
	return isSet
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, isSet := c.items[key]
	if !isSet {
		return nil, false
	}
	c.queue.MoveToFront(item)
	return item.Value.(cachedItem).v, true
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    &sync.Mutex{},
	}
}
