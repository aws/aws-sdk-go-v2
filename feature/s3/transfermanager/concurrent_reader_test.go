package transfermanager

import (
	"bytes"
	"context"
	"io"
	"math"
	"math/rand"
	"sync"
	"testing"
)

func TestConcurrentReader(t *testing.T) {
	cases := map[string]struct {
		partSize     int64
		partsCount   int32
		concurrency  int
		sectionParts int32
	}{
		"single goroutine": {
			partSize:     10,
			partsCount:   1000,
			concurrency:  1,
			sectionParts: 6,
		},
		"single goroutine with only one section": {
			partSize:     1000,
			partsCount:   5,
			concurrency:  3,
			sectionParts: 6,
		},
		"single goroutine with only one part": {
			partSize:     1000,
			partsCount:   1,
			concurrency:  3,
			sectionParts: 6,
		},
		"multiple goroutines": {
			partSize:     10,
			partsCount:   1000,
			concurrency:  5,
			sectionParts: 6,
		},
		"multiple goroutines with only one section": {
			partSize:     10,
			partsCount:   6,
			concurrency:  5,
			sectionParts: 6,
		},
		"multiple goroutines with only one part": {
			partSize:     10,
			partsCount:   1,
			concurrency:  5,
			sectionParts: 6,
		},
		"multiple goroutines with large part size": {
			partSize:     10000,
			partsCount:   10000,
			concurrency:  5,
			sectionParts: 6,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := NewConcurrentReader()
			r.ch = make(chan outChunk, c.concurrency)
			r.setCapacity(int32(math.Min(float64(c.sectionParts), float64(c.partsCount))))
			r.setPartSize(c.partSize)
			r.setPartsCount(c.partsCount)
			ctx := context.Background()
			var wg sync.WaitGroup
			expectBuf := make([]byte, 0)
			actualBuf := make([]byte, 0)

			wg.Add(1)
			go func() {
				defer wg.Done()
				b, err := io.ReadAll(r)
				if err != nil {
					if err != io.EOF {
						t.Error("error copying file: ", err)
					}
					return
				}

				actualBuf = append(actualBuf, b...)
			}()

			getter := mockGetter{}
			ch := make(chan inChunk, c.concurrency)

			for i := 0; i < c.concurrency; i++ {
				getter.wg.Add(1)
				go getter.partGet(ctx, ch, r.ch)
			}

			var i int32
			for {
				if i == c.partsCount {
					break
				}

				if capacity := r.getCapacity(); r.getRead() == capacity {
					r.setCapacity(int32(math.Min(float64(capacity+c.sectionParts), float64(c.partsCount))))
				}

				if i == r.getCapacity() {
					continue
				}

				b := make([]byte, c.partSize)
				if i == c.partsCount-1 {
					b = make([]byte, rand.Intn(int(c.partSize))+1)
				}
				rand.Read(b)
				expectBuf = append(expectBuf, b...)
				ch <- inChunk{
					index: i,
					body:  b,
				}
				i++
			}

			wg.Wait()
			close(ch)
			getter.wg.Wait()
			close(r.ch)

			if e, a := len(expectBuf), len(actualBuf); e != a {
				t.Errorf("expect data sent to have length %d, but got %d", e, a)
			}
			if e, a := expectBuf, actualBuf; !bytes.Equal(e, a) {
				t.Errorf("expect data sent to be %v, got %v", e, a)
			}
		})
	}
}

type mockGetter struct {
	wg sync.WaitGroup
}

func (g *mockGetter) partGet(ctx context.Context, inputCh chan inChunk, outCh chan outChunk) {
	defer g.wg.Done()
	for {
		inC, ok := <-inputCh
		if !ok {
			break
		}

		outCh <- outChunk{
			index:  inC.index,
			body:   bytes.NewReader(inC.body),
			length: int64(len(inC.body)),
		}
	}
}

type inChunk struct {
	body  []byte
	index int32
}
