package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = value

		return true
	}
	c.items[key] = c.queue.PushFront(value)
	lnn := c.queue.Len()
	if lnn > c.capacity {
		toDelItem := c.queue.Back()
		c.queue.Remove(toDelItem)
		for key, val := range c.items { // Удаляем ключ из словаря
			if val == toDelItem {
				delete(c.items, key)
			}
		}
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		newItem := c.queue.Front()
		c.items[key] = newItem
		return newItem.Value, true
	}
	return nil, false
}
