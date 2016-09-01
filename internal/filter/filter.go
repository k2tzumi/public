package filter

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Bloomfilter struct {
	size      uint32
	hashCount int

	mu             sync.Mutex   // mutex to coordinate the transition between bitspaces
	btspSwitchMu   sync.RWMutex // mutext that allows bitspaces to be read while they are updated.
	activeBitspace bool         // which one of the two bitspaces must be read for Has()
	bitspaceF      []int
	bitspaceT      []int
}

func New(size uint32, hashcount int) *Bloomfilter {
	return &Bloomfilter{
		size:      size,
		hashCount: hashcount,
		bitspaceF: make([]int, size),
		bitspaceT: make([]int, size),
	}
}

func (b *Bloomfilter) Add(str string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.btspSwitchMu.Lock()
	defer b.btspSwitchMu.Unlock()

	var bitspaceW []int
	if b.activeBitspace {
		bitspaceW = b.bitspaceF
		copy(bitspaceW, b.bitspaceT)
	} else {
		bitspaceW = b.bitspaceT
		copy(bitspaceW, b.bitspaceF)
	}
	b.activeBitspace = !b.activeBitspace

	for i := 0; i < b.hashCount; i++ {
		idx := b.bitspaceIdx(str, i)
		bitspaceW[idx]++
	}
}

func (b *Bloomfilter) Del(str string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.btspSwitchMu.Lock()
	defer b.btspSwitchMu.Unlock()

	var bitspaceW []int
	if b.activeBitspace {
		bitspaceW = b.bitspaceF
		copy(bitspaceW, b.bitspaceT)
	} else {
		bitspaceW = b.bitspaceT
		copy(bitspaceW, b.bitspaceF)
	}
	b.activeBitspace = !b.activeBitspace

	for i := 0; i < b.hashCount; i++ {
		idx := b.bitspaceIdx(str, i)
		if bitspaceW[idx] > 0 {
			bitspaceW[idx]--
		}
	}
}

func (b *Bloomfilter) Has(str string) bool {
	b.btspSwitchMu.RLock()
	defer b.btspSwitchMu.RUnlock()

	var bitspaceW []int
	if b.activeBitspace {
		bitspaceW = b.bitspaceT
	} else {
		bitspaceW = b.bitspaceF
	}

	for i := 0; i < b.hashCount; i++ {
		idx := b.bitspaceIdx(str, i)
		if bitspaceW[idx] > 0 {
			return true
		}
	}
	return false
}

func (b *Bloomfilter) Saturation() float64 {
	b.btspSwitchMu.RLock()
	defer b.btspSwitchMu.RUnlock()

	var bitspaceW []int
	if b.activeBitspace {
		bitspaceW = b.bitspaceT
	} else {
		bitspaceW = b.bitspaceF
	}

	var i int
	for _, v := range bitspaceW {
		if v > 0 {
			i++
		}
	}
	return float64(i) / float64(len(bitspaceW))
}

func (b *Bloomfilter) bitspaceIdx(str string, i int) uint32 {
	h := fnv.New32()
	fmt.Fprintf(h, "%s%v", str, i)
	sum := h.Sum32()
	return sum % b.size
}
