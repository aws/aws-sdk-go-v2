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

	// the channels we need
	innerErrors := make(chan error, 1)
	outerErrors := make(chan error, 1)

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
			case <-ctx.Done():
				// stop looping
				return
			case <-obj.Lock():
				go func(windowLocation int, obj *windowObject) {
					result, err := producer(windowLocation)
					if err != nil {
						innerErrors <- err
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
		defer close(outChan)
		defer close(innerErrors)
		defer close(outerErrors)

		windowLocation := 0
		for {
			var result interface{}
			obj := window[windowLocation%windowSize]
			select {
			case result = <-obj.out:
				obj.used.Unlock()
			case err := <-innerErrors:
				// any error will stop this
				outerErrors <- err
				return
			case <-ctx.Done():
				return
			}

			// When we get a nil close everything and move on
			if result == nil {
				return
			}

			// Write out the responses
			select {
			case outChan <- result:
			case <-ctx.Done():
				return
			}

			windowLocation += 1
		}
	}()

	// join on the functions
	return outChan, outerErrors
}
