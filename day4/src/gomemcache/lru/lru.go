package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64 // 允许使用的最大内存容量
	nBytes   int64 // 当前已经使用的内存容量
	list     *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len() to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		list:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查询
func (c *Cache) Get(key string) (Value, bool) {
	if element, ok := c.cache[key]; ok {
		// 调整链表中节点的位置
		c.list.MoveToFront(element)
		e := element.Value.(*entry)
		return e.value, true
	}
	return nil, false
}

// 删除
func (c *Cache) RemoveOldest() {
	backElement := c.list.Back()
	if backElement != nil {
		e := backElement.Value.(*entry)
		delete(c.cache, e.key)
		c.list.Remove(backElement)
		// 更新nbytes
		c.nBytes -= int64(len(e.key)) + int64(e.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(e.key, e.value)
		}
	}
}

// 添加
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		e := element.Value.(*entry)
		c.nBytes -= int64(e.value.Len())
		e.value = value
		c.nBytes += int64(value.Len())
		c.list.MoveToFront(element)
	}else {
		element := c.list.PushFront(&entry{key, value})
		c.cache[key] = element
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 如果当前占用内存超过最大可用内存，则移除链表末尾的元素
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// 获取缓存存储的条目数
func (c *Cache) GetEntryCount() int {
	return c.list.Len()
}
