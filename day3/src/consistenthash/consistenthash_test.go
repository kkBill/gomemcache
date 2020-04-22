package consistenthash

import (
	"flag"
	"fmt"
	"strconv"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	// 创建一个一致性哈希实例，并自定义hash函数
	chash := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})
	// 添加真实节点，为了方便说明，这里的节点名称只用数字进行表示
	chash.Add("4", "6", "2")

	testCases := map[string]string{
		"15": "6",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if chash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// 新增一个节点"8"，对应增加3个虚拟节点，分别为8,18,28
	chash.Add("8")

	// 此时如果查询的key为27，将会对应到虚拟节点28，也就是映射到真实节点8
	testCases["27"] = "8"

	for k, v := range testCases {
		if chash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}

var keysPtr = flag.Int("keys", 10000, "key number")
var nodesPtr = flag.Int("nodes", 3, "node number of old cluster")
var newNodesPtr = flag.Int("new-nodes", 4, "node number of new cluster")

// 测试一致性哈希的数据迁移率
func TestMigrateRatio(t *testing.T) {
	flag.Parse()
	var keys = *keysPtr
	var nodes = *nodesPtr
	var newNodes = *newNodesPtr
	fmt.Printf("keys:%d, nodes:%d, newNodes:%d\n", keys, nodes, newNodes)

	c := New(3, nil)
	for i := 0; i < nodes; i++ {
		c.Add(strconv.Itoa(i))
	}

	newC := New(3, nil)
	for i := 0; i < newNodes; i++ {
		newC.Add(strconv.Itoa(i))
	}

	migrate := 0
	for i := 0; i < keys; i++ {
		server := c.Get(strconv.Itoa(i))
		newServer:= newC.Get(strconv.Itoa(i))
		if server != newServer {
			migrate++
		}
	}
	migrateRatio := float64(migrate) / float64(keys)
	fmt.Printf("%f%%\n", migrateRatio*100)
}