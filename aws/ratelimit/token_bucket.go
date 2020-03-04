package ratelimit

import (
	"sync"
)

// TokenBucket provides a concurrency safe utility for adding and removing
// tokens from the available token bucket.
type TokenBucket struct {
	capacity    uint
	maxCapacity uint

	// TODO would it be better to replace this mutex with CAS loops?
	mu sync.Mutex
}

// NewTokenBucket returns an initialized TokenBucket with the capacity
// specified.
func NewTokenBucket(i uint) *TokenBucket {
	return &TokenBucket{
		capacity:    i,
		maxCapacity: i,
	}
}

// Retrieve attempts to reduce the available tokens by the amount requested. If
// available tokens are already negative, or the retrieve would make available
// tokens negative, the retrieve will return false for not retrieved.
// Also returns the available tokens in the bucket.
func (t *TokenBucket) Retrieve(amount uint) (available uint, retrieved bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	a := int(t.capacity - amount)
	if a < 0 {
		return t.capacity, false
	}

	t.capacity = uint(a)
	return t.capacity, true
}

// Refund returns the amount of tokens back to the available token bucket, up
// to the initial capacity.
func (t *TokenBucket) Refund(amount uint) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.capacity += amount
	if t.capacity > t.maxCapacity {
		t.capacity = t.maxCapacity
	}
}
