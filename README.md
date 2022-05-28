# go-sche
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-sche/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-sche/branch/main/graph/badge.svg?token=SQS18BX6SG)](https://codecov.io/gh/czasg/go-sche)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-sche.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-sche/stargazers)

## 背景
任务调度是很常见的需求，一般的任务调度系统包含异步、回调等模块，实现相对比较复杂。  

go-sche 是一个基于 cron、gonal 实现的轻量级任务调度库。

## 目标
1、基于 cron、gonal 实现任务调度
- [x] Cron
- [x] Gonal

2、实现 memory、postgres 存储
- [x] Memory Store
- [x] Postgres Store

## 使用
1、初始化Scheduler
```go
// 依赖
import "github.com/czasg/go-sche"
// 初始化调度
scheduler := sche.NewScheduler()
```

2、新增任务
```go
_ = scheduler.AddTask(&sche.Task{
    Name: "task1",
    Label: map[string]string{
        "gonal标签": "gonal标签",
    },
    Trig: "* * * * *",
})
```

3、启动调度（阻塞）
```go
_ = scheduler.Start(context.Background())
```

4、更新任务（基于ID）
```go
_ = scheduler.AddTask(&sche.Task{
    ID: 1,
    Name: "task1",
    Label: map[string]string{
        "gonal标签": "gonal标签",
    },
    Trig: "*/30 * * * *",
})
```

5、删除任务（基于ID）
```go
_ = scheduler.DelTask(&sche.Task{
    ID: 1,
})
```

## 4.Demo
```go
package main

import (
	"context"
	"fmt"
	"github.com/czasg/go-sche"
	"github.com/czasg/gonal"
	"time"
)

func handler(ctx context.Context, labels gonal.Labels, data []byte) {
	fmt.Println(labels, string(data))
}

func main() {
	// 绑定 gonal 标签
	labels := map[string]string{"test": "test"}
	gonal.BindHandler(labels, handler)
	// 创建任务
	task := sche.Task{
		Name:  "task1",
		Label: labels,
		Trig:  "* * * * *",
	}
	// 初始化调度对象
	scheduler := sche.NewScheduler()
	// 后台挂起调度
	ctx, cancel := context.WithCancel(context.Background())
	go scheduler.Start(ctx)
	// 新增任务
	_ = scheduler.AddTask(&task)
	time.Sleep(time.Second * 5)
	task.Trig = "*/3 * * * *"
	_ = scheduler.UpdateTask(&task)
	time.Sleep(time.Second * 10)
	_ = scheduler.DelTask(&task)
	time.Sleep(time.Second * 5)
	cancel()
}
```
