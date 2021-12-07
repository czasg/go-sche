package sche

import (
	"context"
	"errors"
	"fmt"
	"github.com/czasg/go-queue"
	"github.com/czasg/gonal"
	"testing"
	"time"
)

func TestScheduler_AddTask(t *testing.T) {
	scheduler := Scheduler{}
	assert(t, scheduler.AddTask(&Task{}), fmt.Errorf("Empty spec string"))
	assert(t, scheduler.AddTask(&Task{ID: 1, Trig: "* * * * *"}), StoreInvalidTaskErr)
	assert(t, scheduler.AddTask(&Task{Trig: "* * * * *"}), nil)
}

func TestScheduler_UpdateTask(t *testing.T) {
	scheduler := Scheduler{}
	assert(t, scheduler.UpdateTask(&Task{}), StoreNoTaskErr)
	assert(t, scheduler.AddTask(&Task{Trig: "* * * * *"}), nil)
	assert(t, scheduler.UpdateTask(&Task{ID: 1}), nil)
}

func TestScheduler_DelTask(t *testing.T) {
	scheduler := Scheduler{}
	assert(t, scheduler.DelTask(&Task{}), StoreNoTaskErr)
	assert(t, scheduler.AddTask(&Task{Trig: "* * * * *"}), nil)
	assert(t, scheduler.DelTask(&Task{ID: 1}), nil)
}

func Test_notify(t *testing.T) {
	notifyStoreTodoErr(errors.New("test"))
	notifyStoreUpdateErr(errors.New("test"))
	notifyStoreNextErr(errors.New("test"))
	notifyTaskRunErr(errors.New("test"))
}

func TestWaiter_Wait(t *testing.T) {
	waiter := Waiter{ctx: context.Background()}
	waiter.Wait(time.Now())
	waiter.Wake()
}

var _ Store = (*TestStore)(nil)

type TestStore struct{}

func (t *TestStore) Todo(now time.Time) ([]*Task, error) {
	return nil, errors.New("test")
}

func (t *TestStore) GetNextRunTime() (time.Time, error) {
	panic("implement me")
}

func (t *TestStore) AddTask(task *Task) error {
	panic("implement me")
}

func (t *TestStore) UpdateTask(task *Task) error {
	panic("implement me")
}

func (t *TestStore) DelTask(task *Task) error {
	panic("implement me")
}

func (t *TestStore) GetTaskByID(id int64) (*Task, error) {
	panic("implement me")
}

var _ Store = (*TestStore2)(nil)

type TestStore2 struct{}

func (t *TestStore2) Todo(now time.Time) ([]*Task, error) {
	return []*Task{&Task{}}, nil
}

func (t *TestStore2) GetNextRunTime() (time.Time, error) {
	return MaxDateTime, errors.New("test")
}

func (t *TestStore2) AddTask(task *Task) error {
	panic("implement me")
}

func (t *TestStore2) UpdateTask(task *Task) error {
	return errors.New("test")
}

func (t *TestStore2) DelTask(task *Task) error {
	panic("implement me")
}

func (t *TestStore2) GetTaskByID(id int64) (*Task, error) {
	panic("implement me")
}

func TestScheduler_Start(t *testing.T) {
	scheduler := Scheduler{Store: &TestStore2{}}
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
	scheduler.Start(ctx)

	scheduler = Scheduler{}
	ctx, _ = context.WithTimeout(context.Background(), time.Millisecond)
	scheduler.Start(ctx)
	_ = gonal.SetContext(ctx)
	_ = gonal.SetConcurrent(0)
	_ = gonal.SetQueue(queue.NewFifoMemoryQueue(0))
	//gonal.SetMaxConcurrent(ctx, 0, queue.NewFifoMemoryQueue(0))

	scheduler = Scheduler{}
	ctx, _ = context.WithTimeout(context.Background(), time.Second*2)
	_ = scheduler.AddTask(&Task{Trig: "* * * * *"})
	scheduler.Start(ctx)

	scheduler = Scheduler{Store: &TestStore{}}
	ctx, _ = context.WithTimeout(context.Background(), time.Millisecond)
	scheduler.Start(ctx)
}
