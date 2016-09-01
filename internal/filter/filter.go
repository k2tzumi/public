package filter

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Bloomfilter struct {
	size      uint32
	hashCount int

	mu       sync.Mutex
	bitspace []int
}

func New(size uint32, hashcount int) *Bloomfilter {
	return &Bloomfilter{
		size:      size,
		hashCount: hashcount,
		bitspace:  make([]int, size),
	}
}

func (b *Bloomfilter) Add(str string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := 0; i < b.hashCount; i++ {
		idx := b.bitspaceIdx(str, i)
		b.bitspace[idx]++
	}
}

func (b *Bloomfilter) Del(str string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := 0; i < b.hashCount; i++ {
		idx := b.bitspaceIdx(str, i)
		if b.bitspace[idx] > 0 {
			b.bitspace[idx]--
		}
	}
}

func (b *Bloomfilter) Has(str string) bool {
	btsp := make([]int, b.size)
	b.mu.Lock()
	copy(btsp, b.bitspace)
	b.mu.Unlock()
	for i := 0; i < b.hashCount; i++ {
		idx := b.bitspaceIdx(str, i)
		if btsp[idx] > 0 {
			return true
		}
	}
	return false
}

func (b *Bloomfilter) bitspaceIdx(str string, i int) uint32 {
	h := fnv.New32()
	fmt.Fprintf(h, "%s%v", str, i)
	sum := h.Sum32()
	return sum % b.size
}
