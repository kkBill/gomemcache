package main

import (
	"container/list"
	"fmt"
)

type Entry struct {
	key   int
	value int
	freq  int
}

type Cache struct {
	keyToElem  map[int]*list.Element // 键:key 值:链表节点(节点Value存放的是*Entry)
	freqToList map[int]*list.List    // 键:频率 值:链表(存储访问频率相同的节点)
	capacity   int                   // 缓存允许的最大容量
	size       int                   // 当前存放的条目数
	minFreq    int                   // 记录当前缓存中最小的频率，初始化-1，表示缓存为空
}

// constructor
func New(capacity int) *Cache {
	return &Cache{
		keyToElem:  make(map[int]*list.Element),
		freqToList: make(map[int]*list.List),
		capacity:   capacity,
		size:       0,
		minFreq:    0,
	}
}

func (cache *Cache) Put(key int, value int) {
	if cache.capacity <= 0 {    //特判
		return
	}

	if ele, ok := cache.keyToElem[key]; ok {
		ele.Value.(*Entry).value = value
		cache.IncreaseFrequency(ele)
	} else {
		if cache.size >= cache.capacity {
			ll := cache.freqToList[cache.minFreq]
			ele := ll.Back()
			ll.Remove(ele)
			entry := ele.Value.(*Entry)
			delete(cache.keyToElem, entry.key)
			cache.size--
		}

		entry := &Entry{
			key:   key,
			value: value,
			freq:  1,
		}
		if ll := cache.freqToList[1]; ll == nil {
			cache.freqToList[1] = list.New()
		}
		ele := cache.freqToList[1].PushFront(entry)
		cache.minFreq = 1
		cache.keyToElem[key] = ele
		cache.size++
		//cache.printCache("put ")
	}
}

func (cache *Cache) Get(key int) int {
	if ele, ok := cache.keyToElem[key]; ok {
		cache.IncreaseFrequency(ele)
		return ele.Value.(*Entry).value
	}
	return -1
}

func (cache *Cache) IncreaseFrequency(ele *list.Element) {
	// 从原链表中移除
	entry := ele.Value.(*Entry)
	currFreq := entry.freq
	ll := cache.freqToList[currFreq]
	ll.Remove(ele)
	if ll.Len() == 0 {
		delete(cache.freqToList, currFreq)

		if currFreq == cache.minFreq {
			cache.minFreq++
		}
	}
	// 加入新的链表
	entry.freq++
	ll = cache.freqToList[entry.freq]
	if ll == nil {
		ll = list.New()
		cache.freqToList[entry.freq] = ll
	}
	cache.keyToElem[entry.key] = ll.PushFront(entry)
}

//
func (cache *Cache) printCache(text string) {
	fmt.Printf("****%s****\n", text)
	for freq, ll := range cache.freqToList {
		fmt.Printf("frequency: %d, ", freq)
		for item := ll.Front(); item != nil; item = item.Next() {
			entry := item.Value.(*Entry)
			fmt.Printf("<%d,%d> ", entry.key, entry.value)
		}
		fmt.Println()
	}
	fmt.Printf("************\n")
}

// 测试
func main() {
	//lfu := New(2)
	//lfu.Put(1, 1)
	//lfu.Put(2, 2)
	//fmt.Println(lfu.Get(1)) // 1
	//lfu.Put(3, 3)    // 移除 key 2
	//fmt.Println(lfu.Get(2)) // -1, not found
	//fmt.Println(lfu.Get(3)) // 3
	//lfu.Put(4, 4)    // 移除 key 1
	//fmt.Println(lfu.Get(1)) // -1
	//fmt.Println(lfu.Get(3)) // 3
	//fmt.Println(lfu.Get(4)) // 4
	//lfu.Put(3, 33)
	//fmt.Println(lfu.Get(3)) // 33

	lfu := New(10)
	lfu.Put(10, 13)
	lfu.Put(3, 17)
	lfu.Put(6, 11)
	lfu.Put(10, 5)
	lfu.Put(9, 10)
	fmt.Println(lfu.Get(13)) // -1
	lfu.Put(2, 19)
	fmt.Println(lfu.Get(2))  // 19
	fmt.Println(lfu.Get(3))  // 17
	lfu.Put(5, 25)
	fmt.Println(lfu.Get(8))  // -1
	lfu.Put(9, 22)
	lfu.Put(5, 5)
	lfu.Put(1, 30)
	fmt.Println(lfu.Get(11)) // -1
	lfu.Put(9, 12)
	fmt.Println(lfu.Get(7)) // -1
	fmt.Println(lfu.Get(5)) // 5
	fmt.Println(lfu.Get(8)) // -1
	fmt.Println(lfu.Get(9)) // 12
	lfu.Put(4, 30)
	lfu.Put(9, 3)
	fmt.Println(lfu.Get(9))  // 3
	fmt.Println(lfu.Get(10)) // 5
	fmt.Println(lfu.Get(10)) // 5



}