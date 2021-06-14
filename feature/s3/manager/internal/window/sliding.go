package window

import (
	"context"
	"sync"
)

type windowObject struct {
	used  sync.Mutex
	write sync.Mutex
	out   chan interface{}
}

func (obj *windowObject) Close() {
	obj.write.Lock()
	defer obj.write.Unlock()

	close(obj.out)
	obj.out = nil
}

func (obj *windowObject) Write(i interface{}) {
	obj.write.Lock()
	defer obj.write.Unlock()

	if obj.out != nil {
		obj.out <- i
	}
}

func (obj *windowObject) Lock() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		obj.used.Lock()
		close(c)
	}()
	return c
}

func SlidingWindow(
	ctx context.Context,
	windowSize int,
	producer func(windowLocation int) (interface{}, error),
) (<-chan interface{}, <-chan error) {
	window := make([]*windowObject, windowSize)
	for i := 0; i < windowSize; i++ {
		window[i] = &windowObject{
			out: make(chan interface{}),
		}
	}
	inner, cancel := context.WithCancel(ctx)

	// the channels we need
	errChan := make(chan error, 1)

	// Starting threads
	go func() {
		defer func() {
			for _, obj := range window {
				obj.Close()
			}
		}()

		windowLocation := 0
		// channel
		for {
			obj := window[windowLocation%windowSize]
			select {
			case <-inner.Done():
				// stop looping
				return
			case <-obj.Lock():
				// The go routine per window node
				go func(windowLocation int, obj *windowObject) {
					result, err := producer(windowLocation)
					if err != nil {
						errChan <- err
						cancel()
						return
					}
					obj.Write(result)
				}(windowLocation, obj)
			}
			windowLocation += 1
		}
	}()

	outChan := make(chan interface{})
	// Slide the window
	go func() {
		defer cancel()
		defer close(outChan)
		defer close(errChan)

		windowLocation := 0
		for {
			var result interface{}
			obj := window[windowLocation%windowSize]
			select {
			case result = <-obj.out:
				obj.used.Unlock()
			case <-inner.Done():
				return
			}

			// When we get a nil close everything and move on
			if result == nil {
				return
			}

			// Write out the responses
			select {
			case outChan <- result:
			case <-inner.Done():
				return
			}

			windowLocation += 1
		}
	}()

	// join on the functions
	return outChan, errChan
}
