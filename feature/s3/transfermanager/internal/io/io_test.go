package io

import (
	"sync"
	"testing"
)

func TestBufferPools(t *testing.T) {
	var pools BufferPools

	const size = 64
	pool := pools.Pool(size)
	if pool == nil {
		t.Fatal("Pool returned nil")
	}

	buf, ok := pool.Get().([]byte)
	if !ok {
		t.Fatalf("Pool returned %T, want []byte", pool.Get())
	}
	if cap(buf) != size {
		t.Fatalf("buffer capacity = %d, want %d", cap(buf), size)
	}

	if got := pools.Pool(size); got != pool {
		t.Fatal("Pool returned a different pool for the same size")
	}
	if got := pools.Pool(size + 1); got == pool {
		t.Fatal("Pool returned the same pool for different sizes")
	}
}

func TestBufferPoolsConcurrent(t *testing.T) {
	var pools BufferPools
	const size = 128

	poolsReturned := make(chan *sync.Pool, 32)
	var wg sync.WaitGroup
	for range cap(poolsReturned) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			poolsReturned <- pools.Pool(size)
		}()
	}
	wg.Wait()
	close(poolsReturned)

	var first *sync.Pool
	for pool := range poolsReturned {
		if first == nil {
			first = pool
			continue
		}
		if pool != first {
			t.Fatal("concurrent Pool calls returned different pools")
		}
	}
}
