package gonal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/czasg/go-queue"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var lock sync.RWMutex
var hub *Hub

type Handler func(ctx context.Context, payload Payload)
type Label map[string]string
type Payload struct {
	Label Label
	Body  []byte
}

func Notify(payload Payload) error {
	lock.RLock()
	defer lock.RUnlock()
	return hub.notify(payload)
}

func Bind(handler Handler, selector Label) {
	lock.Lock()
	defer lock.Unlock()
	hub.bind(handler, selector)
}

func SetMaxConcurrent(ctx context.Context, concurrent int, queue queue.Queue) {
	lock.Lock()
	defer lock.Unlock()
	c, cancel := context.WithCancel(ctx)
	newHub := &Hub{
		Ctx:           c,
		CtxCancel:     cancel,
		Concurrent:    concurrent,
		LabelsMatcher: map[string][]Handler{},
		Q:             queue,
		C:             make(chan struct{}, 1),
	}
	if hub != nil {
		newHub.LabelsMatcher = hub.LabelsMatcher
		hub.close()
	}
	hub = newHub
	go hub.loop()
	time.Sleep(time.Millisecond)
}

type Hub struct {
	Ctx           context.Context
	CtxCancel     context.CancelFunc
	Concurrent    int
	LabelsMatcher map[string][]Handler
	Q             queue.Queue
	C             chan struct{}
}

func (h *Hub) notify(payload Payload) error {
	select {
	case <-h.Ctx.Done():
		return h.Ctx.Err()
	default:
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = h.Q.Push(body)
	if err != nil {
		return err
	}
	select {
	case h.C <- struct{}{}:
	default:
	}
	return nil
}

func (h *Hub) bind(handler Handler, selector Label) {
	for k, v := range selector {
		key := fmt.Sprintf(`%s=%v`, k, v)
		pool := h.LabelsMatcher[key]
		if pool == nil {
			pool = []Handler{}
		}
		pool = append(pool, handler)
		h.LabelsMatcher[key] = pool
	}
}

func (h *Hub) fetch(selectors Label) []Handler {
	set := map[reflect.Value]struct{}{}
	handlers := []Handler{}
	for k, v := range selectors {
		key := fmt.Sprintf(`%s=%v`, k, v)
		pool, ok := h.LabelsMatcher[key]
		if !ok {
			continue
		}
		for _, handler := range pool {
			_, ok := set[reflect.ValueOf(handler)]
			if ok {
				continue
			}
			set[reflect.ValueOf(handler)] = struct{}{}
			handlers = append(handlers, handler)
		}
	}
	return handlers
}

func (h *Hub) loop() {
	defer h.close()
	ch := make(chan struct{}, h.Concurrent)
	for {
		select {
		case <-h.Ctx.Done():
			return
		case <-h.C:
		}
		func() {
			for {
				body, err := h.Q.Pop()
				if err != nil {
					notifyQueuePopErr(err)
					return
				}
				var payload Payload
				err = json.Unmarshal(body, &payload)
				if err != nil {
					notifyJsonErr(err)
					return
				}
				for _, handler := range h.fetch(payload.Label) {
					go func(han Handler) {
						defer func() {
							if err := recover(); err != nil {
								notifyHandlerPanic(err)
							}
							select {
							case <-ch:
							default:
							}
						}()
						han(h.Ctx, payload)
					}(handler)
					select {
					case <-h.Ctx.Done():
						return
					case ch <- struct{}{}:
					}
				}
			}
		}()
	}
}

func (h *Hub) close() {
	if h.CtxCancel != nil {
		h.CtxCancel()
	}
}

func notifyQueuePopErr(err error) {
	if errors.Is(err, queue.ErrEmptyQueue) {
		return
	}
	_ = Notify(Payload{
		Label: Label{
			"gonal.internal.event":       "failure",
			"gonal.internal.event.type":  "loop.queue.pop",
			"gonal.internal.event.error": err.Error(),
		},
	})
}

func notifyJsonErr(err error) {
	_ = Notify(Payload{
		Label: Label{
			"gonal.internal.event":       "failure",
			"gonal.internal.event.type":  "loop.json.unmarshal",
			"gonal.internal.event.error": err.Error(),
		},
	})
}

func notifyHandlerPanic(err interface{}) {
	_ = Notify(Payload{
		Label: Label{
			"gonal.internal.event":       "failure",
			"gonal.internal.event.type":  "loop.handler.panic",
			"gonal.internal.event.error": fmt.Sprintf("%v", err),
		},
	})
}

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()
	SetMaxConcurrent(ctx, runtime.NumCPU()*4, queue.NewFifoMemoryQueue(1024))
}
