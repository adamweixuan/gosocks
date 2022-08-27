package main

import (
	"math/bits"
	"sync"
)

const (
	size = 16
)

var (
	caches [size]sync.Pool
)

func init() {
	for i := 0; i < size; i++ {
		bufSize := 1 << i
		caches[i].New = func() any {
			buf := make([]byte, 0, bufSize)
			return buf
		}
	}
}

func getIdx(size uint) int {
	if size == 0 {
		return 0
	}
	if isPowerOfTwo(size) {
		return bsr(size)
	}
	return bsr(size) + 1
}

func isPowerOfTwo(x uint) bool {
	return (x & (-x)) == x
}

func bsr(x uint) int {
	return bits.Len(x) - 1
}

func Get(size uint, capacity uint) []byte {
	c := size
	if capacity > size {
		c = capacity
	}
	var ret = caches[getIdx(c)].Get().([]byte)
	ret = ret[:size]
	return ret
}

func Put(buf []byte) {
	size := uint(cap(buf))
	if !isPowerOfTwo(size) {
		return
	}
	buf = buf[:0]
	caches[bsr(size)].Put(buf)
}
