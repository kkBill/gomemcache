//package bloomfilter
package main

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"github.com/willf/bitset"
	"math"
)

type BloomFilter struct {
	m uint // size of bloom filter
	k uint // number of hash functions
	b *bitset.BitSet
}

func max(x, y uint) uint {
	if x > y {
		return x
	} else {
		return y
	}
}

// 通过哈希函数产生4个基础哈希值
func baseHashes(data []byte) [4]uint64 {
	a1 := []byte{1} // to grab another bit of data
	hasher := murmur3.New128()
	hasher.Write(data) // #nosec
	v1, v2 := hasher.Sum128()
	hasher.Write(a1) // #nosec
	v3, v4 := hasher.Sum128()
	return [4]uint64{
		v1, v2, v3, v4,
	}
}

// location returns the ith hashed location using the four base hash values
func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}

// 通过对哈希值求余，确保位置落在[0,m)范围内
func (bf *BloomFilter) location(h [4]uint64, i uint) uint {
	return uint(location(h, i) % uint64(bf.m))
}

// 根据元素个数n，错误率p，计算出合适的布隆过滤器的参数
// 这里的 Log(2) 就是数学表示中的 ln(2)
// 这里的公式计算推导参考博客，或者参考https://en.wikipedia.org/wiki/Bloom_filter
func EstimateParameters(n uint, p float64) (m uint, k uint) {
	m = uint(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k = uint(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return
}

func New(m, k uint) *BloomFilter {
	return &BloomFilter{
		m: max(m, 1), // 这么做是为了避免错误输入，保证了m至少为正数
		k: max(k, 1), // 同上
		b: bitset.New(m),
	}
}

func NewWithEstimates(n uint, p float64) *BloomFilter {
	m, k := EstimateParameters(n, p)
	return New(m, k)
}

// 对于每个插入的数据项，分别计算出k个哈希值，作为位数组的下标
func (bf *BloomFilter) Add(data []byte) *BloomFilter {
	h := baseHashes(data)
	for i:=uint(0); i < bf.k; i++ {
		bf.b.Set(bf.location(h, i)) //利用第i个哈希值确定对应的需要标记的位置
	}
	return bf
}

// Test returns true if the data is in the BloomFilter, false otherwise.
// If true, the result might be a false positive. If false, the data
// is definitely not in the set.
func (bf *BloomFilter) Test(data []byte)  bool {
	h := baseHashes(data)
	for i := uint(0); i < bf.k; i++ {
		if !bf.b.Test(bf.location(h,i)) { // 如果遇到某一位为0，则必定不存在
			return false
		}
	}
	return true
}

func main() {
	//b := bitset.New(10)
	//b.Set(4)
	//fmt.Println(b.Test(4)) // true
	//fmt.Println(b.Test(5)) // false
	//fmt.Println(b.Test(100)) // false

	filter := New(uint(1024), uint(3))
	filter.Add([]byte("ZJU"))
	fmt.Println(filter.Test([]byte("ZJU"))) // true
	fmt.Println(filter.Test([]byte("ZJNU"))) // false
}