package gonal

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/czasg/go-queue"
    "reflect"
    "runtime"
    "sync"
)

var (
    gonal      *Gonal
    ErrRunning = errors.New("gonal running")
)

func Notify(payload Payload) error {
    return gonal.Notify(payload)
}

func Bind(label Label, handler ...Handler) {
    gonal.Bind(label, handler...)
}

func Fetch(label Label) []Handler {
    return gonal.Fetch(label)
}

func SetContext(ctx context.Context) error {
    return gonal.SetContext(ctx)
}

func SetQueue(q queue.Queue) error {
    return gonal.SetQueue(q)
}

func SetConcurrent(concurrent int) error {
    return gonal.SetConcurrent(concurrent)
}

func Close() error {
    return gonal.Close()
}

type Handler func(ctx context.Context, payload Payload)
type Label map[string]string
type Payload struct {
    Label Label
    Body  []byte
}
type Gonal struct {
    once          sync.Once
    Ctx           context.Context
    Cancel        context.CancelFunc
    Concurrent    int
    LabelsMatcher map[string][]Handler
    Q             queue.Queue
    C             chan struct{}
    Lock          sync.Mutex
    Running       bool
}

func (g *Gonal) Notify(payload Payload) error {
    g.once.Do(func() {
        g.Lock.Lock()
        defer g.Lock.Unlock()
        g.Running = true
        go g.loop()
    })
    select {
    case <-g.Ctx.Done():
        return g.Ctx.Err()
    default:
    }
    body, _ := json.Marshal(payload)
    err := g.Q.Push(body)
    if err != nil {
        return err
    }
    select {
    case g.C <- struct{}{}:
    default:
    }
    return nil
}

func (g *Gonal) Bind(label Label, handlers ...Handler) {
    g.Lock.Lock()
    defer g.Lock.Unlock()
    if len(handlers) < 1 {
        return
    }
    for k, v := range label {
        key := fmt.Sprintf(`%s=%v`, k, v)
        pool := g.LabelsMatcher[key]
        if pool == nil {
            pool = []Handler{}
        }
        pool = append(pool, handlers...)
        g.LabelsMatcher[key] = pool
    }
}

func (g *Gonal) Fetch(label Label) []Handler {
    set := map[reflect.Value]struct{}{}
    handlers := []Handler{}
    for k, v := range label {
        key := fmt.Sprintf(`%s=%v`, k, v)
        pool, ok := g.LabelsMatcher[key]
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

func (g *Gonal) SetContext(ctx context.Context) error {
    g.Lock.Lock()
    defer g.Lock.Unlock()
    if g.Running {
        return ErrRunning
    }
    g.Ctx, g.Cancel = context.WithCancel(ctx)
    return nil
}

func (g *Gonal) SetConcurrent(concurrent int) error {
    g.Lock.Lock()
    defer g.Lock.Unlock()
    if g.Running {
        return ErrRunning
    }
    g.Concurrent = concurrent
    return nil
}

func (g *Gonal) SetQueue(queue queue.Queue) error {
    g.Lock.Lock()
    defer g.Lock.Unlock()
    if g.Running {
        return ErrRunning
    }
    g.Q = queue
    return nil
}

func (g *Gonal) Close() error {
    g.Lock.Lock()
    defer g.Lock.Unlock()
    g.Running = false
    if g.Cancel != nil {
        g.Cancel()
    }
    if g.Q == nil {
        return nil
    }
    return g.Q.Close()
}

func (g *Gonal) loop() {
    defer g.Close()
    ch := make(chan struct{}, g.Concurrent)
    for {
        select {
        case <-g.Ctx.Done():
            return
        case <-g.C:
        }
        func() {
            for {
                body, err := g.Q.Pop()
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
                for _, handler := range g.Fetch(payload.Label) {
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
                        han(g.Ctx, payload)
                    }(handler)
                    select {
                    case <-g.Ctx.Done():
                        return
                    case ch <- struct{}{}:
                    }
                }
            }
        }()
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
    gonal = &Gonal{
        LabelsMatcher: map[string][]Handler{},
        C:             make(chan struct{}, 1),
    }
    _ = SetContext(context.Background())
    _ = SetConcurrent(runtime.NumCPU() * 4)
    _ = SetQueue(queue.NewFifoMemoryQueue(1024))
}
