package io

import (
	"io"
	"sync"
	"unsafe"
)

// Pools is the shared set of buffer pools used by transfer operations.
var Pools BufferPools

// right now we only do direct i/o on linux, according to the google the vast
// majority of storage devices are either 512 or 4k, 4k works with 512 and it's
// not a lot of memory so just align by that
//
// out linux file code does check the block size on init, on the off chance it
// doesn't work w/ 4k we just do normal file i/o
//
// we can make this generic but it's more code than i'd like to own atm without
// evidence that customers are actually using storage devices that need
// arbitrary alignments
const alignedBy = 4096

// BufferPool is a pool of byte buffers. *sync.Pool satisfies this interface.
type BufferPool interface {
	Get() any
	Put(any)
}

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
	Init(size, partSize, writeSize int64, directIO bool) error
	Close() error
}
