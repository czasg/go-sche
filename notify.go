package sche

import (
	"context"
	"time"
)

type Notify struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (n *Notify) Wait(t time.Time) <-chan struct{} {
	ctx, cancel := context.WithDeadline(n.ctx, t)
	n.cancel = cancel
	return ctx.Done()
}

func (n *Notify) Notify() {
	if n.cancel != nil {
		n.cancel()
	}
}
