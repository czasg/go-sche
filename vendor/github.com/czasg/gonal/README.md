# gonal
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/gonal/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/gonal/branch/main/graph/badge.svg?token=XRI6I1W0C3)](https://codecov.io/gh/czasg/gonal)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/gonal.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/gonal/stargazers)

## 1.背景
基于信号的通知/调度，是非常常见的需求，比如在 go 中，就可以通过`signal.Notify`监控系统信号。   
在内部程序调用上，有时为了更好的解耦各个模块之间的依赖，也可以引入信号通知机制。

gonal 提供了信号的**异步通知与处理能力**。

## 2.目标
1、基于标签的多元信号机制
- [x] Labels

2、支持控制并发处理数
- [X] Concurrent

## 3.使用
1、绑定Handler
```go
// 依赖
import "github.com/czasg/gonal"
// handler 函数
func handler(ctx context.Context, labels map[string]string, data []byte) {
    fmt.Println(labels, string(data))
}
// 标签
labels := map[string]string{
    "属性一": "值",
    "属性二": "值",
}

// 绑定标签与函数
gonal.BindHandler(labels, handler)
// 查询绑定函数
_ = gonal.FetchHandler(labels)
```

2、发送通知
```go
// 依赖
import "github.com/czasg/gonal"
// 标签
labels := map[string]string{
    "属性一": "值",
}
// 消息内容
data := []byte("data")

// 发送通知，指定标签与数据
gonal.Notify(nil, labels, data)
```

3、重置最大并发数
```go
gonal.SetConcurrent(1)
```

## 4.Demo
1、绑定与通知
```go
package main

import (
    "context"
    "fmt"
    "github.com/czasg/gonal"
    "time"
)

func handler(ctx context.Context, labels gonal.Labels, data []byte) {
    fmt.Println(labels, data)
}

func main() {
    labels := gonal.Labels{
        "绑定属性一": "一",
        "绑定属性二": "二",
    }
    gonal.BindHandler(labels, handler)

    _ = gonal.Notify(nil, labels, nil)
    time.Sleep(time.Millisecond)
}
```
