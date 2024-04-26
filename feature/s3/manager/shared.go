package manager

type concurrencyLock struct {
	l chan struct{}
}

func newConcurrencyLock(size int) *concurrencyLock {
	return &concurrencyLock{
		l: make(chan struct{}, size),
	}
}

func (c *concurrencyLock) Lock() {
	c.l <- struct{}{}
}

func (c *concurrencyLock) Unlock() {
	<-c.l
}

func (c *concurrencyLock) Close() {
	close(c.l)
}
