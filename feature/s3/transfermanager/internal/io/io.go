package io

import (
	"io"
	"sync"
	"unsafe"
)

var Pools BufferPools

const alignedBy = 4096

// BufferPools retains a separate sync.Pool for each buffer size.
type BufferPools struct {
	mu    sync.Mutex
	pools map[int]*sync.Pool
}

// https://en.wikipedia.org/wiki/Data_structure_alignment#Computing_padding
//
// IMPORTANT: this only works when alignedBy is a power of 2
func align[T uintptr | int64](addr T) T {
	return (addr + (alignedBy - 1)) &^ (alignedBy - 1)
}

func makealigned(size int) []byte {
	p := make([]byte, size+alignedBy)
	u := unsafe.Pointer(&p[0])
	addr := uintptr(u)

	off := align(addr) - addr
	return p[off : int(off)+size : int(off)+size]
}

// Pool returns the pool for buffers of the requested size.
//
// DO NOT repeatedly call Pool() in a transfer operation. Grab the pool of the
// size you need once and retain a reference to it.
func (bps *BufferPools) Pool(size int) *sync.Pool {
	bps.mu.Lock()
	defer bps.mu.Unlock()

	if bps.pools == nil {
		bps.pools = make(map[int]*sync.Pool)
	}
	if p, ok := bps.pools[size]; ok {
		return p
	}

	p := &sync.Pool{
		New: func() any {
			// we only NEED alignment for O_DIRECT on linux but it's only 4k
			// extra bytes so w/e
			return makealigned(size)
		},
	}
	bps.pools[size] = p
	return p
}

// File is a lazily-initialized download destination.
type File interface {
	io.WriterAt
	Init(int64, int64) error
	Close() error
}
