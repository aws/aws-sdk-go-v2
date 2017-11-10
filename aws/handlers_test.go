package aws_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestHandlerList(t *testing.T) {
	s := ""
	r := &aws.Request{}
	l := aws.HandlerList{}
	l.PushBack(func(r *aws.Request) {
		s += "a"
		r.Data = s
	})
	l.Run(r)
	if e, a := "a", s; e != a {
		t.Errorf("expect %q update got %q", e, a)
	}
	if e, a := "a", r.Data.(string); e != a {
		t.Errorf("expect %q data update got %q", e, a)
	}
}

func TestMultipleHandlers(t *testing.T) {
	r := &aws.Request{}
	l := aws.HandlerList{}
	l.PushBack(func(r *aws.Request) { r.Data = nil })
	l.PushFront(func(r *aws.Request) { r.Data = aws.Bool(true) })
	l.Run(r)
	if r.Data != nil {
		t.Error("Expected handler to execute")
	}
}

func TestNamedHandlers(t *testing.T) {
	l := aws.HandlerList{}
	named := aws.NamedHandler{Name: "Name", Fn: func(r *aws.Request) {}}
	named2 := aws.NamedHandler{Name: "NotName", Fn: func(r *aws.Request) {}}
	l.PushBackNamed(named)
	l.PushBackNamed(named)
	l.PushBackNamed(named2)
	l.PushBack(func(r *aws.Request) {})
	if e, a := 4, l.Len(); e != a {
		t.Errorf("expect %d list length, got %d", e, a)
	}
	l.Remove(named)
	if e, a := 2, l.Len(); e != a {
		t.Errorf("expect %d list length, got %d", e, a)
	}
}

func TestSwapHandlers(t *testing.T) {
	firstHandlerCalled := 0
	swappedOutHandlerCalled := 0
	swappedInHandlerCalled := 0

	l := aws.HandlerList{}
	named := aws.NamedHandler{Name: "Name", Fn: func(r *aws.Request) {
		firstHandlerCalled++
	}}
	named2 := aws.NamedHandler{Name: "SwapOutName", Fn: func(r *aws.Request) {
		swappedOutHandlerCalled++
	}}
	l.PushBackNamed(named)
	l.PushBackNamed(named2)
	l.PushBackNamed(named)

	l.SwapNamed(aws.NamedHandler{Name: "SwapOutName", Fn: func(r *aws.Request) {
		swappedInHandlerCalled++
	}})

	l.Run(&aws.Request{})

	if e, a := 2, firstHandlerCalled; e != a {
		t.Errorf("expect first handler to be called %d, was called %d times", e, a)
	}
	if n := swappedOutHandlerCalled; n != 0 {
		t.Errorf("expect swapped out handler to not be called, was called %d times", n)
	}
	if e, a := 1, swappedInHandlerCalled; e != a {
		t.Errorf("expect swapped in handler to be called %d, was called %d times", e, a)
	}
}

func TestSetBackNamed_Exists(t *testing.T) {
	firstHandlerCalled := 0
	swappedOutHandlerCalled := 0
	swappedInHandlerCalled := 0

	l := aws.HandlerList{}
	named := aws.NamedHandler{Name: "Name", Fn: func(r *aws.Request) {
		firstHandlerCalled++
	}}
	named2 := aws.NamedHandler{Name: "SwapOutName", Fn: func(r *aws.Request) {
		swappedOutHandlerCalled++
	}}
	l.PushBackNamed(named)
	l.PushBackNamed(named2)

	l.SetBackNamed(aws.NamedHandler{Name: "SwapOutName", Fn: func(r *aws.Request) {
		swappedInHandlerCalled++
	}})

	l.Run(&aws.Request{})

	if e, a := 1, firstHandlerCalled; e != a {
		t.Errorf("expect first handler to be called %d, was called %d times", e, a)
	}
	if n := swappedOutHandlerCalled; n != 0 {
		t.Errorf("expect swapped out handler to not be called, was called %d times", n)
	}
	if e, a := 1, swappedInHandlerCalled; e != a {
		t.Errorf("expect swapped in handler to be called %d, was called %d times", e, a)
	}
}

func TestSetBackNamed_NotExists(t *testing.T) {
	firstHandlerCalled := 0
	secondHandlerCalled := 0
	swappedInHandlerCalled := 0

	l := aws.HandlerList{}
	named := aws.NamedHandler{Name: "Name", Fn: func(r *aws.Request) {
		firstHandlerCalled++
	}}
	named2 := aws.NamedHandler{Name: "OtherName", Fn: func(r *aws.Request) {
		secondHandlerCalled++
	}}
	l.PushBackNamed(named)
	l.PushBackNamed(named2)

	l.SetBackNamed(aws.NamedHandler{Name: "SwapOutName", Fn: func(r *aws.Request) {
		swappedInHandlerCalled++
	}})

	l.Run(&aws.Request{})

	if e, a := 1, firstHandlerCalled; e != a {
		t.Errorf("expect first handler to be called %d, was called %d times", e, a)
	}
	if e, a := 1, secondHandlerCalled; e != a {
		t.Errorf("expect second handler to be called %d, was called %d times", e, a)
	}
	if e, a := 1, swappedInHandlerCalled; e != a {
		t.Errorf("expect swapped in handler to be called %d, was called %d times", e, a)
	}
}

func TestLoggedHandlers(t *testing.T) {
	expectedHandlers := []string{"name1", "name2"}
	l := aws.HandlerList{}
	loggedHandlers := []string{}
	l.AfterEachFn = aws.HandlerListLogItem
	cfg := aws.Config{Logger: aws.LoggerFunc(func(args ...interface{}) {
		loggedHandlers = append(loggedHandlers, args[2].(string))
	})}

	named1 := aws.NamedHandler{Name: "name1", Fn: func(r *aws.Request) {}}
	named2 := aws.NamedHandler{Name: "name2", Fn: func(r *aws.Request) {}}
	l.PushBackNamed(named1)
	l.PushBackNamed(named2)
	l.Run(&aws.Request{Config: cfg})

	if !reflect.DeepEqual(expectedHandlers, loggedHandlers) {
		t.Errorf("expect handlers executed %v to match logged handlers, %v",
			expectedHandlers, loggedHandlers)
	}
}

func TestStopHandlers(t *testing.T) {
	l := aws.HandlerList{}
	stopAt := 1
	l.AfterEachFn = func(item aws.HandlerListRunItem) bool {
		return item.Index != stopAt
	}

	called := 0
	l.PushBackNamed(aws.NamedHandler{Name: "name1", Fn: func(r *aws.Request) {
		called++
	}})
	l.PushBackNamed(aws.NamedHandler{Name: "name2", Fn: func(r *aws.Request) {
		called++
	}})
	l.PushBackNamed(aws.NamedHandler{Name: "name3", Fn: func(r *aws.Request) {
		t.Fatalf("third handler should not be called")
	}})
	l.Run(&aws.Request{})

	if e, a := 2, called; e != a {
		t.Errorf("expect %d handlers called, got %d", e, a)
	}
}

func BenchmarkNewRequest(b *testing.B) {
	svc := s3.New(unit.Config())

	for i := 0; i < b.N; i++ {
		r := svc.GetObjectRequest(nil)
		if r.Request == nil {
			b.Fatal("r should not be nil")
		}
	}
}

func BenchmarkHandlersCopy(b *testing.B) {
	handlers := aws.Handlers{}

	handlers.Validate.PushBack(func(r *aws.Request) {})
	handlers.Validate.PushBack(func(r *aws.Request) {})
	handlers.Build.PushBack(func(r *aws.Request) {})
	handlers.Build.PushBack(func(r *aws.Request) {})
	handlers.Send.PushBack(func(r *aws.Request) {})
	handlers.Send.PushBack(func(r *aws.Request) {})
	handlers.Unmarshal.PushBack(func(r *aws.Request) {})
	handlers.Unmarshal.PushBack(func(r *aws.Request) {})

	for i := 0; i < b.N; i++ {
		h := handlers.Copy()
		if e, a := handlers.Validate.Len(), h.Validate.Len(); e != a {
			b.Fatalf("expected %d handlers got %d", e, a)
		}
	}
}

func BenchmarkHandlersPushBack(b *testing.B) {
	handlers := aws.Handlers{}

	for i := 0; i < b.N; i++ {
		h := handlers.Copy()
		h.Validate.PushBack(func(r *aws.Request) {})
		h.Validate.PushBack(func(r *aws.Request) {})
		h.Validate.PushBack(func(r *aws.Request) {})
		h.Validate.PushBack(func(r *aws.Request) {})
	}
}

func BenchmarkHandlersPushFront(b *testing.B) {
	handlers := aws.Handlers{}

	for i := 0; i < b.N; i++ {
		h := handlers.Copy()
		h.Validate.PushFront(func(r *aws.Request) {})
		h.Validate.PushFront(func(r *aws.Request) {})
		h.Validate.PushFront(func(r *aws.Request) {})
		h.Validate.PushFront(func(r *aws.Request) {})
	}
}

func BenchmarkHandlersClear(b *testing.B) {
	handlers := aws.Handlers{}

	for i := 0; i < b.N; i++ {
		h := handlers.Copy()
		h.Validate.PushFront(func(r *aws.Request) {})
		h.Validate.PushFront(func(r *aws.Request) {})
		h.Validate.PushFront(func(r *aws.Request) {})
		h.Validate.PushFront(func(r *aws.Request) {})
		h.Clear()
	}
}
