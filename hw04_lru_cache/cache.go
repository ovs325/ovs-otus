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

type ItemVal struct {
	KeyCache Key
	ValItem  any
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = ItemVal{KeyCache: key, ValItem: value}
		return true
	}

	if c.queue.Len() == c.capacity {
		delete(c.items, c.queue.Back().Value.(ItemVal).KeyCache)
		c.queue.Remove(c.queue.Back())
	}
	c.items[key] = c.queue.PushFront(ItemVal{KeyCache: key, ValItem: value})
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.Value.(ItemVal).ValItem, true
	}
	return nil, false
}
