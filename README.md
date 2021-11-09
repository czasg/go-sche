# go-sche
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-sche/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-sche/branch/main/graph/badge.svg?token=SQS18BX6SG)](https://codecov.io/gh/czasg/go-sche)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-sche.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-sche/stargazers)

go-sche is a golang library that lets you schedule your task to be executed later.

You can add new jobs or remove old ones on the fly as you please. 
If you store your jobs in a database like postgres,
they will also survive scheduler restarts and maintain their state. 

When the scheduler restarted, it will then run all the jobs it should have run while it was offline.
```text
|—————————————|                          notify by gonal
|  scheduler  | ————————> task<labels> |-----------------> label |-----> handler
|—————————————|                        |                         |-----> handler
       |  interface                    |                         |-----> handler
|—————————————|                        |          
|    store    |                        |-----------------> label |-----> handler
|—————————————|                        |                         | ...
       |  next run time                | ...
|—————————————|
|    task     |
|—————————————|
       |  last run time
|—————————————|
|  cron-trig  |
|—————————————|
```

# demo
```go
package main

import (
	"context"
	"fmt"
	"github.com/czasg/go-sche"
	"github.com/czasg/gonal"
	"time"
)

func worker1(ctx context.Context, payload gonal.Payload) {
	fmt.Println("worker1", payload.Label, time.Now())
}

func worker2(ctx context.Context, payload gonal.Payload) {
	fmt.Println("worker2", payload.Label, time.Now())
}

func init() {
	// bind task with labels by gonal.
	gonal.Bind(worker1, gonal.Label{"func": "worker1"})
	gonal.Bind(worker2, gonal.Label{"func": "worker2"})
}

func main() {
	// init
	scheduler := sche.Scheduler{}
	// add task
	_ = scheduler.AddTask(&sche.Task{
		Name: "task1",
		Trig: "*/20 * * * *",
		Label: map[string]string{
			"func":       "worker1",
			"meta.other": "test1",
		},
	})
	_ = scheduler.AddTask(&sche.Task{
		Name: "task2",
		Trig: "15 * * * *",
		Label: map[string]string{
			"func":       "worker2",
			"meta.other": "test2",
		},
	})
	// start with block.
	scheduler.Start(context.Background())
}
```

# more
### using postgres
```go
package main

import (
	"context"
	"github.com/czasg/go-sche"
)

func main() {
	pg := NewPG()
	
	scheduler := sche.Scheduler{
		Store: sche.NewPostgresStore(pg),
	}
	scheduler.Start(context.Background())
}
```
