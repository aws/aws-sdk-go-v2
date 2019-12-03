package sdkio

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	byteSlice := make([]byte, 10)
	cases := map[string]struct {
		expectedStart int
		expectedEnd   int
		expectedSize  int
	}{
		"ringBuff": {},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			ringBuffer := NewRingBuffer(byteSlice)
			if e, o := c.expectedSize, ringBuffer.size; e != o {
				t.Errorf("expect default size to be %v , got %v", e, o)
			}
			if e, o := c.expectedStart, ringBuffer.start; e != o {
				t.Errorf("expect deafult start to point to %v , got %v", e, o)
			}
			if e, o := c.expectedEnd, ringBuffer.end; e != o {
				t.Errorf("expect default end tp point to %v , got %v", e, o)
			}

		})
	}
}

func TestRingBuffer_Write(t *testing.T) {
	cases := map[string]struct {
		sliceCapacity         int
		input                 []byte
		expectedStart         int
		expectedEnd           int
		expectedSize          int
		expectedWrittenBuffer []byte
	}{
		"Capacity matches Bytes written": {
			sliceCapacity:         11,
			input:                 []byte("hello world"),
			expectedStart:         0,
			expectedEnd:           11,
			expectedSize:          11,
			expectedWrittenBuffer: []byte("hello world"),
		},
		"Capacity is lower than Bytes written": {
			sliceCapacity:         10,
			input:                 []byte("hello world"),
			expectedStart:         1,
			expectedEnd:           1,
			expectedSize:          10,
			expectedWrittenBuffer: []byte("dello worl"),
		},
		"Capacity is more than Bytes written": {
			sliceCapacity:         12,
			input:                 []byte("hello world"),
			expectedStart:         0,
			expectedEnd:           11,
			expectedSize:          11,
			expectedWrittenBuffer: []byte("hello world"),
		},
		"No Bytes written": {
			sliceCapacity:         10,
			input:                 []byte(""),
			expectedStart:         0,
			expectedEnd:           0,
			expectedSize:          0,
			expectedWrittenBuffer: []byte(""),
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			byteSlice := make([]byte, c.sliceCapacity)
			ringBuffer := NewRingBuffer(byteSlice)
			ringBuffer.Write(c.input)

			if e, a := c.expectedSize, ringBuffer.size; e != a {
				t.Errorf("expect default size to be %v , got %v", e, a)
			}

			if e, a := c.expectedStart, ringBuffer.start; e != a {
				t.Errorf("expect deafult start to point to %v , got %v", e, a)
			}

			if e, a := c.expectedEnd, ringBuffer.end; e != a {
				t.Errorf("expect default end to point to %v , got %v", e, a)
			}

			if e, a := c.expectedWrittenBuffer, ringBuffer.slice; !bytes.Contains(a, e) {
				t.Errorf("expect written bytes to be %v , got %v", string(e), string(a))
			}
		})
	}

}

func TestRingBuffer_Read(t *testing.T) {
	cases := map[string]struct {
		input                         []byte
		numberOfBytesToRead           int
		expectedStartAfterRead        int
		expectedEndAfterRead          int
		expectedSizeOfBufferAfterRead int
		expectedReadSlice             []byte
	}{
		"Read capacity matches Bytes written": {
			input:                         []byte("Hello world"),
			numberOfBytesToRead:           11,
			expectedStartAfterRead:        0,
			expectedEndAfterRead:          11,
			expectedSizeOfBufferAfterRead: 0,
			expectedReadSlice:             []byte("Hello world"),
		},
		"Read capacity is lower than Bytes written": {
			input:                         []byte("hello world"),
			numberOfBytesToRead:           5,
			expectedStartAfterRead:        5,
			expectedEndAfterRead:          11,
			expectedSizeOfBufferAfterRead: 6,
			expectedReadSlice:             []byte("hello"),
		},
		"Read capacity is more than Bytes written": {
			input:                         []byte("hello world"),
			numberOfBytesToRead:           15,
			expectedStartAfterRead:        0,
			expectedEndAfterRead:          11,
			expectedSizeOfBufferAfterRead: 0,
			expectedReadSlice:             []byte("hello world"),
		},
		"No Bytes are read": {
			input:                         []byte("hello world"),
			numberOfBytesToRead:           0,
			expectedStartAfterRead:        0,
			expectedEndAfterRead:          11,
			expectedSizeOfBufferAfterRead: 11,
			expectedReadSlice:             []byte(""),
		},
		"No Bytes written": {
			input:                         []byte(""),
			numberOfBytesToRead:           11,
			expectedStartAfterRead:        0,
			expectedEndAfterRead:          0,
			expectedSizeOfBufferAfterRead: 0,
			expectedReadSlice:             []byte(""),
		},
		"Write capacity is more than Bytes Written": {
			input:                         []byte("h"),
			numberOfBytesToRead:           11,
			expectedStartAfterRead:        1,
			expectedEndAfterRead:          1,
			expectedSizeOfBufferAfterRead: 0,
			expectedReadSlice:             []byte("h"),
		},
	}
	for _, c := range cases {
		byteSlice := make([]byte, 11)
		t.Run("", func(t *testing.T) {
			ringBuffer := NewRingBuffer(byteSlice)
			readSlice := make([]byte, c.numberOfBytesToRead)

			ringBuffer.Write(c.input)
			_, err := ringBuffer.Read(readSlice)

			if e, a := io.EOF, err; e != a {
				t.Errorf("Expected %v, got %v", e, a)
			}
			if e, a := c.expectedReadSlice, readSlice; !bytes.Contains(a, e) {
				t.Errorf("expect read buffer to be %v, got %v", string(e), string(a))
			}

			if e, a := c.expectedSizeOfBufferAfterRead, ringBuffer.size; e != a {
				t.Errorf("expect default size to be %v , got %v", e, a)
			}

			if e, a := c.expectedStartAfterRead, ringBuffer.start; e != a {
				t.Errorf("expect default start to point to %v , got %v", e, a)
			}

			if e, a := c.expectedEndAfterRead, ringBuffer.end; e != a {
				t.Errorf("expect default end to point to %v , got %v", e, a)
			}
		})
	}
}

func TestRingBuffer_forConsecutiveReadWrites(t *testing.T) {
	cases := map[string]struct {
		input                         []string
		sliceCapacity                 int
		numberOfBytesToRead           []int
		expectedStartAfterRead        []int
		expectedEnd                   []int
		expectedSizeOfBufferAfterRead []int
		expectedReadSlice             []string
		expectedWrittenBuffer         []string
	}{
		"Capacity matches Bytes written": {
			input:                         []string{"Hello World", "Hello Earth", "Mars, "},
			sliceCapacity:                 11,
			numberOfBytesToRead:           []int{5, 11},
			expectedStartAfterRead:        []int{5, 6},
			expectedEnd:                   []int{11, 6},
			expectedSizeOfBufferAfterRead: []int{6, 0},
			expectedReadSlice:             []string{"Hello", "EarthMars, "},
			expectedWrittenBuffer:         []string{"Hello World", "Hello Earth", "Mars, Earth"},
		},
		// "Capacity is lower than Bytes written": {},
		// "Capacity is more than Bytes written":  {},
		// "No Bytes written":                     {},
	}
	for _, c := range cases {
		writeSlice := make([]byte, c.sliceCapacity)
		ringBuffer := NewRingBuffer(writeSlice)

		t.Run("", func(t *testing.T) {

			// write: 0
			ringBuffer.Write([]byte(c.input[0]))
			if e, a := string(ringBuffer.slice), c.expectedWrittenBuffer[0]; e != a {
				t.Errorf("Expected %v, got %v", e, a)
			}

			// read: 0
			readSlice := make([]byte, c.numberOfBytesToRead[0])
			readCount, err := ringBuffer.Read(readSlice)

			if e, a := io.EOF, err; e != a {
				t.Errorf("Expected %v, got %v", e, a)
			}

			if e, a := len(c.expectedReadSlice[0]), readCount; e != a {
				t.Errorf("Expected to read %v bytes, read only %v", e, a)
			}

			if e, a := c.expectedReadSlice[0], string(readSlice); !strings.Contains(a, e) {
				t.Errorf("expect read buffer to be %v, got %v", e, a)
			}

			if e, a := c.expectedSizeOfBufferAfterRead[0], ringBuffer.size; e != a {
				t.Errorf("expect buffer size to be %v , got %v", e, a)
			}

			if e, a := c.expectedStartAfterRead[0], ringBuffer.start; e != a {
				t.Errorf("expect default start to point to %v , got %v", e, a)
			}

			if e, a := c.expectedEnd[0], ringBuffer.end; e != a {
				t.Errorf("expect default end tp point to %v , got %v", e, a)
			}

			/*
				Next cycle of read writes.
			*/

			// write: 1
			ringBuffer.Write([]byte(c.input[1]))
			if e, a := c.expectedWrittenBuffer[1], string(ringBuffer.slice); e != a {
				t.Errorf("Expected %v, got %v", e, a)
			}

			// write: 2
			ringBuffer.Write([]byte(c.input[2]))
			if e, a := c.expectedWrittenBuffer[2], string(ringBuffer.slice); e != a {
				t.Errorf("Expected %v, got %v", e, a)
			}

			// read: 1
			readSlice = make([]byte, c.numberOfBytesToRead[1])
			readCount, err = ringBuffer.Read(readSlice)

			if e, a := io.EOF, err; e != a {
				t.Errorf("Expected %v, got %v", e, a)
			}

			if e, a := len(c.expectedReadSlice[1]), readCount; e != a {
				t.Errorf("Expected to read %v bytes, read only %v", e, a)
			}

			if e, a := c.expectedReadSlice[1], string(readSlice); e != a {
				t.Errorf("expect read buffer to be %v, got %v", e, a)
			}

			if e, a := c.expectedSizeOfBufferAfterRead[1], ringBuffer.size; e != a {
				t.Errorf("expect buffer size to be %v , got %v", e, a)
			}

			if e, a := c.expectedStartAfterRead[1], ringBuffer.start; e != a {
				t.Errorf("expect default start to point to %v , got %v", e, a)
			}

			if e, a := c.expectedEnd[1], ringBuffer.end; e != a {
				t.Errorf("expect default end to point to %v , got %v", e, a)
			}

		})
	}
}
