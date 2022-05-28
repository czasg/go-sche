package sche

import (
	"context"
	"testing"
	"time"
)

func TestScheduler_AddTask(t *testing.T) {
	scheduler := Scheduler{}
	assert(t, scheduler.AddTask(&Task{}), nil)
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

func TestScheduler_Start(t *testing.T) {
	label := map[string]string{"test": "1"}
	scheduler := Scheduler{}
	err := scheduler.AddTask(&Task{
		ID:          0,
		Name:        "",
		Label:       label,
		Trig:        "* * * * *",
		LastRunTime: time.Time{},
		NextRunTime: time.Time{},
		Suspended:   false,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go scheduler.Start(ctx)
	time.Sleep(time.Second * 2)
	cancel()
}
