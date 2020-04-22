package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 哈希函数签名
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // 哈希函数
	replicas int            // 每个真实节点对应的虚拟节点个数
	circle   []int          // 哈希环
	hashMap  map[int]string // 存储虚拟节点与真实节点的映射
}

func New(replicas int, hashFunc Hash) *Map {
	m := &Map{
		hash:     hashFunc,
		replicas: replicas,
		hashMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 添加真实节点
func (m *Map) Add(peers ...string) {
	for _, peer := range peers {
		for i := 0; i < m.replicas; i++ {
			// 根据节点的名称+编号计算hash值
			hashVal := int(m.hash([]byte(strconv.Itoa(i) + peer)))
			//log.Printf("peer:[%s-%d], hash:[%d]\n", peer, i, hashVal)
			m.circle = append(m.circle, hashVal)
			m.hashMap[hashVal] = peer
		}
	}
	// 排序是为了接下来方便查询
	sort.Ints(m.circle)
	//log.Printf("circle:\n %v\n",m.circle)
}

// 根据请求的key选择对应的真实节点
// 返回真实节点的名称
func (m *Map) Get(key string) string {
	if len(m.circle) == 0 {
		return ""
	}
	// 首先计算key对应的hash值
	hashVal := int(m.hash([]byte(key)))
	// 二分查找
	index := sort.Search(len(m.circle), func(i int) bool {
		return m.circle[i] >= hashVal
	})
	return m.hashMap[m.circle[index%len(m.circle)]]
}
