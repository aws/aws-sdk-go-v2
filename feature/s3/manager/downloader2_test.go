package manager

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
)

func Test_receiveChunks(t *testing.T) {
	tests := []struct {
		input       []int
		concurrency int
	}{
		{
			input:       []int{0, 4, 1, 2, 3, 5, 7, 6},
			concurrency: 5,
		},
		{
			input:       []int{0, 2, 6, 5, 4, 1, 3, 10, 9, 7, 8},
			concurrency: 6,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v", tt.input), func(t *testing.T) {
			ch := make(chan *dlchunk2, tt.concurrency)
			var g sync.WaitGroup
			g.Add(1)
			go func() {
				defer g.Done()
				for _, i := range tt.input {
					ch <- testNewChunk(i)
				}
				close(ch)
			}()

			var b bytes.Buffer
			err := receiveChunks(&b, tt.concurrency, ch)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			g.Wait()

			var builder strings.Builder
			for seq := 0; seq < len(tt.input); seq++ {
				builder.WriteString(fmt.Sprintf("hello there %d\n", seq))
			}

			if e, a := builder.String(), b.String(); e != a {
				t.Errorf("expect %q response, got %q", e, a)
			}
		})
	}

}

func testNewChunk(seq int) *dlchunk2 {
	part := make([]byte, 20)
	chunk := &dlchunk2{
		part:    &part,
		cleanup: func() {},
		seq:     seq,
		size:    int64(len(part)),
	}
	_, _ = io.WriteString(chunk, fmt.Sprintf("hello there %d\n", seq))
	return chunk
}
