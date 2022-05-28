package gonal

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/czasg/go-queue"
    "reflect"
    "runtime"
    "sync"
    "time"
)

func NewGonal(ctx context.Context, q queue.Queue) *Gonal {
    if ctx == nil {
        ctx = context.Background()
    }
    if q == nil {
        q = queue.NewFifoMemoryQueue()
    }
    gonal := &Gonal{
        queue:    q,
        handlers: map[string][]Handler{},
    }
    gonal.ctx, gonal.cancel = context.WithCancel(ctx)
    gonal.SetConcurrent(runtime.NumCPU() * 2)
    return gonal
}

func Notify(ctx context.Context, labels Labels, data []byte) error {
    return gonal.Notify(ctx, labels, data)
}

func BindHandler(labels Labels, handlers ...Handler) {
    gonal.BindHandler(labels, handlers...)
}

func FetchHandler(labels Labels) []Handler {
    return gonal.FetchHandler(labels)
}

func SetConcurrent(concurrent int) {
    gonal.SetConcurrent(concurrent)
}

func Close() error {
    return gonal.Close()
}

var gonal = NewGonal(nil, nil)

type Placeholder struct{}
type Labels map[string]string
type Handler func(ctx context.Context, labels Labels, data []byte)

type Gonal struct {
    ctx      context.Context
    cancel   context.CancelFunc
    reset    context.CancelFunc
    handlers map[string][]Handler
    queue    queue.Queue
    lock     sync.Mutex
}

type Payload struct {
    Labels Labels
    Data   []byte
}

func (g *Gonal) Notify(ctx context.Context, labels Labels, data []byte) error {
    select {
    case <-g.ctx.Done():
        return g.ctx.Err()
    default:
    }
    body, err := json.Marshal(Payload{Labels: labels, Data: data})
    if err != nil {
        return err
    }
    return g.queue.Put(ctx, body)
}

func (g *Gonal) BindHandler(labels Labels, handlers ...Handler) {
    g.lock.Lock()
    defer g.lock.Unlock()
    for key, value := range labels {
        hk := fmt.Sprintf("%s=%s", key, value)
        g.handlers[hk] = append(g.handlers[hk], handlers...)
    }
}

func (g *Gonal) FetchHandler(labels Labels) []Handler {
    g.lock.Lock()
    defer g.lock.Unlock()
    results := []Handler{}
    handlerSet := map[reflect.Value]Placeholder{}
    for key, value := range labels {
        hk := fmt.Sprintf("%s=%s", key, value)
        handlers, ok := g.handlers[hk]
        if !ok {
            continue
        }
        for _, handler := range handlers {
            _, ok := handlerSet[reflect.ValueOf(handler)]
            if ok {
                continue
            }
            handlerSet[reflect.ValueOf(handler)] = Placeholder{}
            results = append(results, handler)
        }
    }
    return results
}

func (g *Gonal) Close() error {
    if g.cancel != nil {
        g.cancel()
    }
    return g.queue.Close()
}

func (g *Gonal) SetConcurrent(concurrent int) {
    g.lock.Lock()
    defer g.lock.Unlock()
    if g.reset != nil {
        g.reset()
        time.Sleep(time.Millisecond * 100)
    }
    ctx, cancel := context.WithCancel(g.ctx)
    limit := make(chan Placeholder, concurrent)
    g.reset = cancel
    go g.loop(ctx, limit)
}

func (g *Gonal) loop(ctx context.Context, limit chan Placeholder) {
    defer close(limit)
    for {
        data, err := g.queue.Get(ctx)
        if err != nil {
            return
        }
        var payload Payload
        err = json.Unmarshal(data, &payload)
        if err != nil {
            continue
        }
        for _, handler := range g.FetchHandler(payload.Labels) {
            select {
            case <-ctx.Done():
                return
            case limit <- Placeholder{}:
            }
            go func(handler Handler) {
                defer func() {
                    if err := recover(); err != nil {}
                    <-limit
                }()
                handler(ctx, payload.Labels, payload.Data)
            }(handler)
        }
    }
}
