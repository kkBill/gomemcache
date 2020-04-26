package main

import (
	"flag"
	"fmt"
	"gomemcache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 创建缓存实例
func createGroup() *gomemcache.Group {
	return gomemcache.NewGroup("scores", 1024, gomemcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Printf("search key[%s] from db.", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("key[%s] does not exist.", key)
		}))
}

// 启动缓存服务器
func startCacheServer(addr string, addrs []string, gocache *gomemcache.Group) {
	server := gomemcache.NewHTTPPool(addr)
	server.Set(addrs...)
	gocache.RegisterPeers(server)
	log.Println("gomemcache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr[len("http://"):], server))
}

//
func startAPIServer(apiAddr string, gocache *gomemcache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gocache.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[len("http://"):], nil))
}

func main() {
	var port int // 指定端口启动HTTP服务
	var api bool
	flag.IntVar(&port, "port", 8001, "gocache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gocache := createGroup()
	if api {
		go startAPIServer(apiAddr, gocache)
	}
	startCacheServer(addrMap[port], addrs, gocache)
}
