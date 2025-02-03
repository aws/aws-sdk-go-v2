package transfermanager

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

// these tests also cover the cache impl (cycling+expiry+evict)

type mockNow struct {
	now time.Time
}

func (m *mockNow) Now() time.Time {
	return m.now
}

func (m *mockNow) Add(d time.Duration) {
	m.now = m.now.Add(d)
}

func useMockNow(m *mockNow) func() {
	timeNow = m.Now
	return func() {
		timeNow = time.Now
	}
}

var errDialContextOK = errors.New("dial context ok")

type mockResolver struct {
	addrs map[string][]string
	err   error
}

func (m *mockResolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	return m.addrs[host], m.err
}

type mockDialContext struct {
	calledWith string
}

func (m *mockDialContext) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	m.calledWith = addr
	return nil, errDialContextOK
}

func TestRoundRobinDNS_CycleIPs(t *testing.T) {
	restore := useMockNow(&mockNow{})
	defer restore()

	addrs := []string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}
	r := &mockResolver{
		addrs: map[string][]string{
			"s3.us-east-1.amazonaws.com": addrs,
		},
	}
	dc := &mockDialContext{}

	rr := &rrDNS{
		cache:       newDNSCache(1),
		resolver:    r,
		dialContext: dc.DialContext,
	}

	expectDialContext(t, rr, dc, "s3.us-east-1.amazonaws.com", addrs[0])
	expectDialContext(t, rr, dc, "s3.us-east-1.amazonaws.com", addrs[1])
	expectDialContext(t, rr, dc, "s3.us-east-1.amazonaws.com", addrs[2])
	expectDialContext(t, rr, dc, "s3.us-east-1.amazonaws.com", addrs[0])
}

func TestRoundRobinDNS_MultiIP(t *testing.T) {
	restore := useMockNow(&mockNow{})
	defer restore()

	r := &mockResolver{
		addrs: map[string][]string{
			"host1.com": {"0.0.0.1", "0.0.0.2", "0.0.0.3"},
			"host2.com": {"1.0.0.1", "1.0.0.2", "1.0.0.3"},
		},
	}
	dc := &mockDialContext{}

	rr := &rrDNS{
		cache:       newDNSCache(2),
		resolver:    r,
		dialContext: dc.DialContext,
	}

	expectDialContext(t, rr, dc, "host1.com", r.addrs["host1.com"][0])
	expectDialContext(t, rr, dc, "host2.com", r.addrs["host2.com"][0])
	expectDialContext(t, rr, dc, "host1.com", r.addrs["host1.com"][1])
	expectDialContext(t, rr, dc, "host2.com", r.addrs["host2.com"][1])
}

func TestRoundRobinDNS_MaxHosts(t *testing.T) {
	restore := useMockNow(&mockNow{})
	defer restore()

	r := &mockResolver{
		addrs: map[string][]string{
			"host1.com": {"0.0.0.1", "0.0.0.2", "0.0.0.3"},
			"host2.com": {"0.0.0.1", "0.0.0.2", "0.0.0.3"},
		},
	}
	dc := &mockDialContext{}

	rr := &rrDNS{
		cache:       newDNSCache(1),
		resolver:    r,
		dialContext: dc.DialContext,
	}

	expectDialContext(t, rr, dc, "host1.com", r.addrs["host1.com"][0])
	expectDialContext(t, rr, dc, "host1.com", r.addrs["host1.com"][1])
	expectDialContext(t, rr, dc, "host2.com", r.addrs["host2.com"][0]) // evicts host1
	expectDialContext(t, rr, dc, "host1.com", r.addrs["host1.com"][0]) // evicts host2
	expectDialContext(t, rr, dc, "host2.com", r.addrs["host2.com"][0])
}

func TestRoundRobinDNS_Expires(t *testing.T) {
	now := &mockNow{time.Unix(0, 0)}
	restore := useMockNow(now)
	defer restore()

	r := &mockResolver{
		addrs: map[string][]string{
			"host1.com": {"0.0.0.1", "0.0.0.2", "0.0.0.3"},
		},
	}
	dc := &mockDialContext{}

	rr := &rrDNS{
		cache:       newDNSCache(2),
		expiry:      30,
		resolver:    r,
		dialContext: dc.DialContext,
	}

	expectDialContext(t, rr, dc, "host1.com", r.addrs["host1.com"][0])
	now.Add(16) // hasn't expired
	expectDialContext(t, rr, dc, "host1.com", r.addrs["host1.com"][1])
	now.Add(16) // expired, starts over
	expectDialContext(t, rr, dc, "host1.com", r.addrs["host1.com"][0])
}

func expectDialContext(t *testing.T, rr *rrDNS, dc *mockDialContext, host, expect string) {
	const port = "443"

	t.Helper()
	_, err := rr.DialContext(context.Background(), "", net.JoinHostPort(host, port))
	if err != errDialContextOK {
		t.Errorf("expect sentinel err, got %v", err)
	}
	actual, _, err := net.SplitHostPort(dc.calledWith)
	if err != nil {
		t.Fatal(err)
	}
	if expect != actual {
		t.Errorf("expect addr %s, got %s", expect, actual)
	}
}
