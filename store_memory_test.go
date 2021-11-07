package sche

import (
	"reflect"
	"testing"
	"time"
)

func assert(t *testing.T, v1, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		t.Error(v1, v2)
	}
}

func TestNewStoreMemory(t *testing.T) {
	now := time.Now()
	store := NewStoreMemory()
	assert(t, store.AddTask(&Task{ID: 1}), StoreInvalidTaskErr)
	assert(t, store.AddTask(&Task{Name: "task1", NextRunTime: now.Add(time.Second)}), nil)
	assert(t, store.AddTask(&Task{Name: "task2", NextRunTime: now.Add(time.Minute)}), nil)
	assert(t, store.AddTask(&Task{Name: "task3", NextRunTime: now.Add(time.Hour)}), nil)
	assert(t, store.UpdateTask(&Task{ID: 3, Name: "task3", NextRunTime: now.Add(time.Hour)}), nil)
	assert(t, store.AddTask(&Task{Name: "task4", NextRunTime: now.Add(time.Minute * 2)}), nil)
	assert(t, store.UpdateTask(&Task{}), StoreNoTaskErr)
	assert(t, store.UpdateTask(&Task{ID: 1, NextRunTime: time.Now().Add(time.Second)}), nil)
	assert(t, store.DelTask(&Task{ID: 3}), nil)
	assert(t, store.DelTask(&Task{}), StoreNoTaskErr)
	task, err := store.GetTaskByID(1)
	assert(t, task, &Task{ID: 1, NextRunTime: now.Add(time.Second)})
	assert(t, err, nil)
	_, err = store.GetTaskByID(0)
	assert(t, err, StoreNoTaskErr)
	next, err := store.GetNextRunTime()
	assert(t, next.Unix(), time.Now().Add(time.Second).Unix())
	assert(t, err, nil)

	store = NewStoreMemory()
	next, err = store.GetNextRunTime()
	assert(t, next, MaxDateTime)
	assert(t, err, nil)
	assert(t, store.AddTask(&Task{Name: "task1", Suspended: true}), nil)
	next, err = store.GetNextRunTime()
	assert(t, next, MaxDateTime)
	assert(t, err, nil)

	store = NewStoreMemory()
	_ = store.AddTask(&Task{Suspended: true})
	_ = store.AddTask(&Task{NextRunTime: now.Add(-time.Second)})
	_ = store.AddTask(&Task{NextRunTime: now.Add(time.Second)})
	_, _ = store.Todo(now)
}
