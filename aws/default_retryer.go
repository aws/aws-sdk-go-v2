package aws

import (
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// DefaultRetryer implements basic retry logic using exponential backoff for
// most services. If you want to implement custom retry logic, implement the
// Retryer interface or create a structure type that composes this
// struct and override the specific methods. For example, to override only
// the MaxRetries method:
//
//		type retryer struct {
//      client.DefaultRetryer
//    }
//
//    // This implementation always has 100 max retries
//    func (d retryer) MaxRetries() int { return 100 }
type DefaultRetryer struct {
	NumMaxRetries    int
	MinRetryDelay    time.Duration
	MinThrottleDelay time.Duration
	MaxRetryDelay    time.Duration
	MaxThrottleDelay time.Duration
}

const (
	// DefaultRetryerMaxNumRetries sets max number of retries
	DefaultRetryerMaxNumRetries = 3

	// DefaultRetryerMinRetryDelay sets min retry delay
	DefaultRetryerMinRetryDelay = 30 * time.Millisecond

	// DefaultRetryerMinThrottleDelay sets minimum delay when throttled
	DefaultRetryerMinThrottleDelay = 500 * time.Millisecond

	// DefaultRetryerMaxRetryDelay sets max retry delay
	DefaultRetryerMaxRetryDelay = 300 * time.Second

	// DefaultRetryerMaxThrottleDelay sets maximum delay when throttled
	DefaultRetryerMaxThrottleDelay = 300 * time.Second
)

// MaxRetries returns the number of maximum returns the service will use to make
// an individual API
func (d DefaultRetryer) MaxRetries() int {
	return d.NumMaxRetries
}

var seededRand = rand.New(&lockedSource{src: rand.NewSource(time.Now().UnixNano())})

// setDefaults sets default values for the default Retryer
func (d *DefaultRetryer) setDefaults() {
	if d.NumMaxRetries == 0 {
		d.NumMaxRetries = DefaultRetryerMaxNumRetries
	}
	if d.MinRetryDelay == 0 {
		d.MinRetryDelay = DefaultRetryerMinRetryDelay
	}
	if d.MinThrottleDelay == 0 {
		d.MinThrottleDelay = DefaultRetryerMinThrottleDelay
	}
	if d.MaxRetryDelay == 0 {
		d.MaxRetryDelay = DefaultRetryerMaxRetryDelay
	}
	if d.MaxThrottleDelay == 0 {
		d.MaxThrottleDelay = DefaultRetryerMaxThrottleDelay
	}

}

// RetryRules returns the delay duration before retrying this request again
// Note: RetryRules method must be a value receiver so that the
// defaultRetryer is safe.
func (d DefaultRetryer) RetryRules(r *Request) time.Duration {

	// set default values for the retryer if not set.
	d.setDefaults()

	// Set the upper limit of delay in retrying at ~five minutes
	minDelay := d.MinRetryDelay
	var initialDelay time.Duration
	throttle := d.shouldThrottle(r)
	if throttle {
		if delay, ok := getRetryAfterDelay(r); ok {
			initialDelay = delay
		}
		minDelay = d.MinThrottleDelay
	}

	retryCount := r.RetryCount
	maxDelay := d.MaxRetryDelay
	if throttle {
		maxDelay = d.MaxThrottleDelay
	}

	var delay time.Duration

	// Logic to cap the retry count based on the minDelay provided
	actualRetryCount := int(math.Log2(float64(minDelay))) + 1
	if actualRetryCount < 63-retryCount {
		delay = time.Duration(1<<uint64(retryCount)) * getDelaySeed(minDelay)
		if delay > maxDelay {
			delay = getDelaySeed(maxDelay / 2)
		}
	} else {
		delay = getDelaySeed(maxDelay / 2)
	}
	return delay + initialDelay
}

func getDelaySeed(duration time.Duration) time.Duration {
	return time.Duration(seededRand.Int63n(int64(duration)) + int64(duration))
}

// ShouldRetry returns true if the request should be retried.
func (d DefaultRetryer) ShouldRetry(r *Request) bool {
	// If one of the other handlers already set the retry state
	// we don't want to override it based on the service's state
	if r.Retryable != nil {
		return *r.Retryable
	}

	if r.HTTPResponse.StatusCode >= 500 {
		return true
	}
	return r.IsErrorRetryable() || d.shouldThrottle(r)
}

// ShouldThrottle returns true if the request should be throttled.
func (d DefaultRetryer) shouldThrottle(r *Request) bool {
	if r.HTTPResponse != nil {
		switch r.HTTPResponse.StatusCode {
		case 429:
		case 502:
		case 503:
		case 504:
		default:
			return r.IsErrorThrottle()
		}
		return true
	}
	return r.IsErrorThrottle()
}

// This will look in the Retry-After header, RFC 7231, for how long
// it will wait before attempting another request
func getRetryAfterDelay(r *Request) (time.Duration, bool) {
	if !canUseRetryAfterHeader(r) {
		return 0, false
	}

	delayStr := r.HTTPResponse.Header.Get("Retry-After")
	if len(delayStr) == 0 {
		return 0, false
	}

	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		return 0, false
	}

	return time.Duration(delay) * time.Second, true
}

// Will look at the status code to see if the retry header pertains to
// the status code.
func canUseRetryAfterHeader(r *Request) bool {
	switch r.HTTPResponse.StatusCode {
	case 429:
	case 503:
	default:
		return false
	}

	return true
}

// lockedSource is a thread-safe implementation of rand.Source
type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}
