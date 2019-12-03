package sdkio

import (
	"io"
)

// RingBuffer struct satisfies io.ReadWrite
// interface.
//
// ReadBuffer is a revolving buffer data structure
// which can be used to store snapshots of data in a
// revolving window.
type RingBuffer struct {
	slice []byte
	start int
	end   int
	size  int
}

// NewRingBuffer method takes in the capacity as an int
// and returns a RingBuffer.
func NewRingBuffer(slice []byte) *RingBuffer {
	ringBuf := RingBuffer{
		slice: slice,
	}
	return &ringBuf
}

// Write method inserts the elements in a byte slice,
// Returns the number of bytes written and an error
//
// This satisfies io.Writer interface.
func (r *RingBuffer) Write(p []byte) (int, error) {

	for _, b := range p {

		// check if end points to invalid index
		// we need to circle back
		if r.end == len(r.slice) {
			r.end = 0
		}

		// check if start points to invalid index
		// we need to circle back
		if r.start == len(r.slice) {
			r.start = 0
		}

		// if ring buffer is filled,
		// increment the start index
		if r.size == len(r.slice) {
			r.size--
			r.start++
		}

		r.slice[r.end] = b
		r.end++
		r.size++
	}
	return r.size, nil
}

// Read method on RingBuffer returns the read count
// along with Error encountered while reading.
//
// Read copies the data on the ring buffer into
// the byte slice provided to the method.
func (r *RingBuffer) Read(p []byte) (int, error) {

	// if ring buffer is empty
	// return EOF error.
	if r.size == 0 {
		return 0, io.EOF
	}

	// readCount keeps track of the number of
	// bytes read.
	readCount := 0

	s := r.start
	for j := 0; j < len(p); j++ {
		p[j] = r.slice[s]
		s++
		readCount++
		r.start++
		r.size--

		if r.start == len(r.slice) {
			r.start = 0
		}

		if s == r.end {
			break
		}

		if s == len(r.slice) {
			s = 0
		}
	}

	return readCount, io.EOF
}
