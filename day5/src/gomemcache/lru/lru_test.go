package lru

import (
	"fmt"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lruCache := New(1024, nil)
	lruCache.Add("key1", String("1234"))
	if value, ok := lruCache.Get("key1"); !ok || value.(String) != String("1234") {
		t.Fatal("cache hit key1=1234 failed")
	}
	if _, ok := lruCache.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"

	capacity := len(k1 + k2 + v1 + v2); // 可用最大内存
	lruCache := New(int64(capacity), nil)
	lruCache.Add(k1, String(v1)) // 当前链表：(k1,v1)
	lruCache.Add(k2, String(v2)) // 当前链表：(k2,v2) -> (k1,v1)
	lruCache.Add(k3, String(v3)) // 当前链表：(k3,v3) -> (k2,v2)，由于内存超限，需要移除(k1, v1)

	if _, ok := lruCache.Get(k1); ok || lruCache.GetEntryCount() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestAdd(t *testing.T) {
	// nothing
}

func TestOnEvicted(t *testing.T) {
	callback := func(key string, value Value) {
		fmt.Printf("entry:<%v, %v> has been deleted.\n",key,value)
	}
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"

	capacity := len(k1 + k2 + v1 + v2); // 可用最大内存
	lruCache := New(int64(capacity), callback)
	lruCache.Add(k1, String(v1)) // 当前链表：(k1,v1)
	lruCache.Add(k2, String(v2)) // 当前链表：(k2,v2) -> (k1,v1)
	lruCache.Add(k3, String(v3)) // 当前链表：(k3,v3) -> (k2,v2)，由于内存超限，需要移除(k1, v1)

	if _, ok := lruCache.Get(k1); ok || lruCache.GetEntryCount() != 2{
		t.Fatalf("Removeoldest key1 failed")
	}
}
