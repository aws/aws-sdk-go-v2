package transfermanager

import (
	"context"
	"sync"
	"testing"
)

func TestDefaultSlicePool(t *testing.T) {
	pool := newDefaultSlicePool(1, 2)

	var bs []byte
	var err error
	var wg sync.WaitGroup

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool.Put(make([]byte, 1))
		}()
	}
	// wait for a slice to be put back
	for i := 0; i < 200; i++ {
		bs, err = pool.Get(context.Background())
		if err != nil {
			t.Errorf("failed to get slice from pool: %v", err)
		}
	}

	wg.Wait()

	// failed to get a slice due to ctx cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	bs, err = pool.Get(ctx)
	if err == nil {
		pool.Put(bs)
		t.Errorf("expectd no slice to be returned")
	}

	if e, a := 2, len(pool.slices); e != a {
		t.Errorf("expect pool size to be %v, got %v", e, a)
	}

	pool.Close()
}
