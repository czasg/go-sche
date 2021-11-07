# gonal
gonal is a set of label, registered-machine, label-selector due to a system.

it's very like signal, but gonal use label to replace it. 

# demo
```golang
package main

import (
	"context"
	"fmt"
	"github.com/czasg/gonal"
	"time"
)

func worker1(ctx context.Context, payload gonal.Payload) { fmt.Println("worker1", payload.Label) }

func worker2(ctx context.Context, payload gonal.Payload) { fmt.Println("worker2", payload.Label) }

func main() {
	gonal.Bind(worker1, gonal.Label{"func": "worker1"})
	gonal.Bind(worker1, gonal.Label{"type": "function"})
	gonal.Bind(worker1, gonal.Label{"meta": "main.worker1"})
	gonal.Bind(worker2, gonal.Label{
		"func": "worker2",
		"type": "function",
		"meta": "main.worker2",
	})

	index := 0
	for {
		for _, label := range []gonal.Label{
			{"func": "worker1"},  // selector will match <worker1>
			{"func": "worker2"},  // selector will match <worker2>
			{"type": "function"}, // selector will match <worker1>&<worker2>
		} {
			index++
			label["index"] = fmt.Sprintf("%05d", index)
			_ = gonal.Notify(gonal.Payload{
				Label: label,
			})
			time.Sleep(time.Second * 5)
		}
	}
}
```

# more
### SetMaxConcurrent
set the max concurrent num. eg:
```golang
package main

import (
	"context"
	"github.com/czasg/go-queue"
	"github.com/czasg/gonal"
	"io/ioutil"
	"os"
)

func main() {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	q, _ := queue.NewFifoDiskQueue(dir)
	gonal.SetMaxConcurrent(context.Background(), 0, q)
}
```
such like concurrent is 0, all payload will serial execution.
### Bind
bind handler with labels.
### Notify
push payload into queue, then match handler by label-selector. 