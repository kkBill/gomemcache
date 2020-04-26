package gomemcache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T)  {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed.")
	}
}

// 利用map简单的模拟数据流
var db = map[string]string {
	"Tom":"630",
	"Jack":"589",
	"Sam":"567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))

	// 创建缓存实例，并指定回调函数
	cache := NewGroup("scores", 1024, GetterFunc(
		func(key string) ([]byte, error) {
			log.Printf("search key[%s] from db.\n", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("key[%s] does not exist.\n", key)
		}))

	// 模拟查询
	for k, v := range db {
		// 第1次查询k，缓存未命中，因此需要从数据库中获取数据
		if res, err := cache.Get(k); err != nil || res.String() != v {
			t.Errorf("fail to get value of [%s] in db.\n", k)
		}
		// 第2次查询k，正常情况下应该从缓存中获取到对应的值
		if _, err := cache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Errorf("key[%s] missed in cache.\n", k)
		}
	}

	if res, err := cache.Get("unknown"); err == nil {
		t.Errorf("the value of key[unknown] should be empty, but [%s] got.\n", res)
	}
}