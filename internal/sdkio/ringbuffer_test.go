package sdkio

import (
	"bytes"
	"io"
	"testing"
)

func TestNew(t *testing.T) {

	cases := []struct {
		name          string
		capacity      int
		expectedStart int
		expectedEnd   int
		expectedSize  int
	}{
		{
			name:          "ringBuff",
			capacity:      10,
			expectedSize:  0,
			expectedStart: 0,
			expectedEnd:   0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ringBuffer := New(c.capacity)

			if e, o := c.capacity, len(ringBuffer.slice); e != o {
				t.Errorf("expect capacity of the ring buffer to be %v , got %v", e, o)
			}

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

	cases := []struct {
		name                  string
		capacity              int
		input                 []byte
		expectedStart         int
		expectedEnd           int
		expectedSize          int
		expectedWrittenBuffer []byte
	}{
		{
			name:                  "Capacity matches Bytes written",
			capacity:              11,
			input:                 []byte("Hello world"),
			expectedStart:         0,
			expectedEnd:           11,
			expectedSize:          11,
			expectedWrittenBuffer: []byte("Hello world"),
		},
		{
			name:                  "Capacity is lower than Bytes written",
			capacity:              9,
			input:                 []byte("hello world"),
			expectedStart:         2,
			expectedEnd:           2,
			expectedSize:          9,
			expectedWrittenBuffer: []byte("ldllo wor"),
		},
		{
			name:                  "Capacity is more than Bytes written",
			capacity:              15,
			input:                 []byte("hello world"),
			expectedStart:         0,
			expectedEnd:           11,
			expectedSize:          11,
			expectedWrittenBuffer: []byte("hello world"),
		},
		{
			name:                  "No Bytes written",
			capacity:              11,
			input:                 []byte(""),
			expectedStart:         0,
			expectedEnd:           0,
			expectedSize:          0,
			expectedWrittenBuffer: []byte(""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ringBuffer := New(c.capacity)
			ringBuffer.Write(c.input)

			if e, a := c.capacity, len(ringBuffer.slice); e != a {
				t.Errorf("expect capacity of the ring buffer to be %v , got %v", e, a)
			}

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
				t.Errorf("expect written bytes to be %v , got %v", e, a)
			}
		})
	}

}

func TestRingBuffer_Read(t *testing.T) {
	cases := []struct {
		name                  string
		capacity              int
		input                 []byte
		sizeofReadBytes       int
		expectedStart         int
		expectedEnd           int
		expectedSize          int
		expectedWrittenBuffer []byte
		expectedReadByteSlice []byte
	}{
		{
			name:                  "Capacity matches Bytes written",
			capacity:              11,
			input:                 []byte("Hello world"),
			expectedStart:         0,
			expectedEnd:           11,
			expectedSize:          11,
			sizeofReadBytes:       11,
			expectedReadByteSlice: []byte("Hello world"),
		},
		{
			name:                  "Capacity is lower than Bytes written",
			capacity:              9,
			input:                 []byte("hello world"),
			expectedStart:         2,
			expectedEnd:           2,
			expectedSize:          9,
			sizeofReadBytes:       9,
			expectedReadByteSlice: []byte("llo world"),
		},
		{
			name:                  "Capacity is more than Bytes written",
			capacity:              15,
			input:                 []byte("hello world"),
			expectedStart:         0,
			expectedEnd:           11,
			expectedSize:          11,
			sizeofReadBytes:       15,
			expectedReadByteSlice: []byte("hello world"),
		},
		{
			name:                  "No Bytes written",
			capacity:              11,
			input:                 []byte(""),
			expectedStart:         0,
			expectedEnd:           0,
			expectedSize:          0,
			sizeofReadBytes:       0,
			expectedReadByteSlice: []byte(""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ringBuffer := New(c.capacity)
			readSlice := make([]byte, c.sizeofReadBytes)

			ringBuffer.Write(c.input)
			_, err := ringBuffer.Read(readSlice)

			if e, a := io.EOF, err; e != a {
				t.Errorf("Expected %v, got %v", e, a)
			}
			if e, a := c.expectedReadByteSlice, readSlice; !bytes.Contains(a, e) {
				t.Errorf("expect read buffer to be %v, got %v", e, a)
			}

			if e, a := c.capacity, len(ringBuffer.slice); e != a {
				t.Errorf("expect capacity of the ring buffer to be %v , got %v", e, a)
			}

			if e, a := c.expectedSize, ringBuffer.size; e != a {
				t.Errorf("expect default size to be %v , got %v", e, a)
			}

			if e, a := c.expectedStart, ringBuffer.start; e != a {
				t.Errorf("expect deafult start to point to %v , got %v", e, a)
			}

			if e, a := c.expectedEnd, ringBuffer.end; e != a {
				t.Errorf("expect default end tp point to %v , got %v", e, a)
			}
		})
	}

}
