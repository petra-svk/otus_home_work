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

type cacheItem struct {
	key   Key
	value interface{}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	newValue := cacheItem{key, value}
	listItem, exist := lc.items[key]

	if exist {
		listItem.Value = newValue
		lc.queue.MoveToFront(listItem)
	} else {
		if lc.capacity == lc.queue.Len() {
			lastItemValue := lc.queue.Back().Value
			lc.queue.Remove(lc.queue.Back())
			delete(lc.items, lastItemValue.(cacheItem).key)
		}
		newListItem := lc.queue.PushFront(newValue)
		lc.items[key] = newListItem
	}
	return exist
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	listItem, exist := lc.items[key]

	if !exist {
		return nil, false
	}

	val := listItem.Value.(cacheItem).value
	lc.queue.MoveToFront(listItem)
	return val, true
}

func (lc *lruCache) Clear() {
	lc.queue.Clear()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
