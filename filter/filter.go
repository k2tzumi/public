package filter

import (
	"fmt"
	"hash/fnv"
)

type bloomfilter struct {
	size      uint32
	hashCount int
	bitspace  []bool
}

func newBloomfilter(size uint32, hashcount int) *bloomfilter {
	return &bloomfilter{
		size:      size,
		hashCount: hashcount,
		bitspace:  make([]bool, size),
	}
}

func (b *bloomfilter) add(str string) {
	for i := 0; i < b.hashCount; i++ {
		idx := b.bitspaceIdx(str, i)
		b.bitspace[idx] = true
	}
}

func (b *bloomfilter) has(str string) bool {
	for i := 0; i < b.hashCount; i++ {
		idx := b.bitspaceIdx(str, i)
		if b.bitspace[idx] {
			return true
		}
	}
	return false
}

func (b *bloomfilter) bitspaceIdx(str string, i int) uint32 {
	h := fnv.New32()
	fmt.Fprintf(h, "%s%v", str, i)
	sum := h.Sum32()
	return sum % b.size
}
