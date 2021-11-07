package sche

import (
	"context"
	"github.com/czasg/gonal"
	"github.com/robfig/cron"
	"time"
)

type Scheduler struct {
	Store
	Waiter
}

func (s *Scheduler) Start(ctx context.Context) {
	var wait time.Time
	s.Waiter = Waiter{ctx: ctx}
	if s.Store == nil {
		s.Store = NewStoreMemory()
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.Waiter.Wait(wait):
		}
		now := time.Now()
		todos, err := s.Store.Todo(now)
		if err != nil {
			wait = MaxDateTime
			notifyStoreTodoErr(err)
			continue
		}
		for _, todo := range todos {
			err = todo.Run()
			if err != nil {
				notifyTaskRunErr(err)
				continue
			}
			err = s.Store.UpdateTask(todo)
			if err != nil {
				notifyStoreUpdateErr(err)
			}
		}
		wait, err = s.Store.GetNextRunTime()
		if err != nil {
			wait = MaxDateTime
			notifyStoreNextErr(err)
		}
	}
}

func (s *Scheduler) AddTask(task *Task) error {
	ins, err := cron.Parse(string(task.Trig))
	if err != nil {
		return err
	}
	TrigCronPool[task.Trig] = ins
	if task.Label == nil {
		task.Label = map[string]string{}
	}
	task.Label["task.label.name"] = task.Name
	task.NextRunTime = task.Trig.GetNextRunTime(time.Now())
	if s.Store == nil {
		s.Store = NewStoreMemory()
	}
	err = s.Store.AddTask(task)
	if err != nil {
		return err
	}
	s.Waiter.Wake()
	return nil
}

func (s *Scheduler) UpdateTask(task *Task) error {
	if task.Label == nil {
		task.Label = map[string]string{}
	}
	task.Label["task.label.name"] = task.Name
	task.NextRunTime = task.Trig.GetNextRunTime(time.Now())
	if s.Store == nil {
		s.Store = NewStoreMemory()
	}
	err := s.Store.UpdateTask(task)
	if err != nil {
		return err
	}
	s.Waiter.Wake()
	return nil
}

func (s *Scheduler) DelTask(task *Task) error {
	if s.Store == nil {
		s.Store = NewStoreMemory()
	}
	err := s.Store.DelTask(task)
	if err != nil {
		return err
	}
	s.Waiter.Wake()
	return nil
}

type Waiter struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *Waiter) Wake() {
	if w.cancel != nil {
		w.cancel()
	}
}

func (w *Waiter) Wait(wait time.Time) <-chan struct{} {
	ctx, cancel := context.WithDeadline(w.ctx, wait)
	w.cancel = cancel
	return ctx.Done()
}

func notifyStoreTodoErr(err error) {
	_ = gonal.Notify(gonal.Payload{Label: gonal.Label{
		"sche.internal.event":       "failure",
		"sche.internal.event.type":  "task.todo",
		"sche.internal.event.error": err.Error(),
	}})
}

func notifyStoreUpdateErr(err error) {
	_ = gonal.Notify(gonal.Payload{Label: gonal.Label{
		"sche.internal.event":       "failure",
		"sche.internal.event.type":  "task.update",
		"sche.internal.event.error": err.Error(),
	}})
}

func notifyStoreNextErr(err error) {
	_ = gonal.Notify(gonal.Payload{Label: gonal.Label{
		"sche.internal.event":       "failure",
		"sche.internal.event.type":  "task.next",
		"sche.internal.event.error": err.Error(),
	}})
}

func notifyTaskRunErr(err error) {
	_ = gonal.Notify(gonal.Payload{Label: gonal.Label{
		"sche.internal.event":       "failure",
		"sche.internal.event.type":  "task.run",
		"sche.internal.event.error": err.Error(),
	}})
}
