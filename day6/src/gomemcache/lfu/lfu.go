package main

//import (
//	"container/list"
//	"fmt"
//)
//
//type CacheItem struct {
//	key             string        // key of entry
//	value           interface{}   // value of item
//	frequencyParent *list.Element // pointer to parent in cacheList
//}
//
//type FrequencyItem struct {
//	entries map[*CacheItem]byte // set of entries
//	freq    int                 // access frequency
//}
//
//type Cache struct {
//	bykey    map[string]*CacheItem // hashmap
//	freqs    *list.List            // doubled-linked list of freq
//	capacity int                   // max numbers of items
//	size     int                   // current size of cache
//}
//
//// constructor
//func New(capacity int) *Cache {
//	return &Cache{
//		bykey:    make(map[string]*CacheItem),
//		freqs:    list.New(),
//		capacity: capacity,
//		size:     0,
//	}
//}
//
//func (cache *Cache) Put(key string, value interface{}) {
//	if item, ok := cache.bykey[key]; ok {
//		item.value = value
//		cache.increment(item)
//	}else {
//		item = &CacheItem{
//			key:             key,
//			value:           value,
//			frequencyParent: nil,
//		}
//		cache.bykey[key] = item
//		cache.size++
//		// Increment item access frequency
//		cache.increment(item)
//		// Eviction, if needed
//		if cache.size > cache.capacity {
//			cache.Evict(1)
//		}
//	}
//}
//
//func (cache *Cache) Get(key string) interface{} {
//	if item, ok := cache.bykey[key]; ok {
//		cache.increment(item)
//		return item.value
//	}
//	return nil
//}
//
//func (cache *Cache) increment(item *CacheItem) {
//	currentFrequency := item.frequencyParent
//	var nextFrequencyAmount int
//	var nextFrequency *list.Element
//
//	// 对于新增条目，其frequencyParent初始为nil
//	if currentFrequency == nil {
//		nextFrequencyAmount = 1
//		nextFrequency = cache.freqs.Front()
//	}else {
//		nextFrequencyAmount = currentFrequency.Value.(*FrequencyItem).freq + 1
//		nextFrequency = currentFrequency.Next()
//	}
//
//	// 需要新建一个frequency 节点
//	if nextFrequency == nil || nextFrequency.Value.(*FrequencyItem).freq != nextFrequencyAmount {
//		newFrequencyItem := &FrequencyItem{
//			entries: make(map[*CacheItem]byte),
//			freq:    nextFrequencyAmount,
//		}
//		// frequency 为1的节点
//		if currentFrequency == nil {
//			nextFrequency = cache.freqs.PushFront(newFrequencyItem)
//		}else {
//			nextFrequency = cache.freqs.InsertAfter(newFrequencyItem, currentFrequency)
//		}
//	}
//
//	item.frequencyParent = nextFrequency
//	nextFrequency.Value.(*FrequencyItem).entries[item] = 1
//	// 对于非新增节点，需要把条目从原先的frequency item对应的列表中移除
//	if currentFrequency != nil {
//		cache.Remove(currentFrequency, item)
//	}
//}
//
//func (cache *Cache) Remove(freqItem *list.Element, cacheItem *CacheItem) {
//	frequencyItem := freqItem.Value.(*FrequencyItem)
//	delete(frequencyItem.entries, cacheItem)
//	if len(frequencyItem.entries) == 0 {
//		cache.freqs.Remove(freqItem)
//	}
//}
//
//// 移除count个条目
//func (cache *Cache) Evict(count int) {
//	for i := 0; i < count; {
//		if freqItem := cache.freqs.Front(); freqItem != nil {
//			for cacheItem, _ := range freqItem.Value.(*FrequencyItem).entries {
//				if i < count {
//					delete(cache.bykey, cacheItem.key)
//					cache.Remove(freqItem, cacheItem)
//					cache.size--
//					i++;
//				}
//			}
//		}
//	}
//}
//
//// 测试
//func main() {
//	lfu := New(2)
//	lfu.Put("1","1")
//	lfu.Put("2","2")
//	fmt.Println(lfu.Get("1")) // 1
//	lfu.Put("3", "3")
//	fmt.Println(lfu.Get("2")) // nil
//	fmt.Println(lfu.Get("3")) // 3
//	lfu.Put("4", "4")
//	fmt.Println(lfu.Get("1")) // 1
//	fmt.Println(lfu.Get("3")) // 3
//	fmt.Println(lfu.Get("4")) // nil
//	lfu.Put("3","33")
//	fmt.Println(lfu.Get("3")) // 33
//}