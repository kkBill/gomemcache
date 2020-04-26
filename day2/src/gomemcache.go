package gomemcache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//Group可以认为是缓存命名空间
type Group struct {
	name   string
	getter Getter
	cache  cache
}

var	(
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:   name,
		getter: getter,
		cache:  cache{maxBytes:maxBytes},
	}
	groups[name] = group
	return group
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	group := groups[name]
	mu.RUnlock()
	return group
}

// Get value for a key from cache in group g
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required.\n")
	}
	if v, ok := g.cache.get(key); ok {
		log.Printf("Group:[%s], Key:[%s], cache hit successfully.\n", g.name, key)
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

// 从本地执行的数据源中加载数据
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	// 一定要注意这一步，在缓存中，存储的是从数据库中读出来数据的副本
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.cache.add(key, value)
}