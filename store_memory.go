package sche

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

var _ Store = (*StoreMemory)(nil)

func NewStoreMemory() *StoreMemory {
	return &StoreMemory{
		Index:    0,
		Tasks:    list.New(),
		TasksMap: map[int64]*list.Element{},
		Lock:     sync.RWMutex{},
	}
}

type StoreMemory struct {
	Index    int64
	Tasks    *list.List
	TasksMap map[int64]*list.Element
	Lock     sync.RWMutex
}

func (s *StoreMemory) Todo(now time.Time) ([]*Task, error) {
	tasks := []*Task{}
	for el := s.Tasks.Front(); el != nil; el = el.Next() {
		task := el.Value.(*Task)
		if task.Suspended {
			continue
		}
		if task.NextRunTime.Before(now) {
			tasks = append(tasks, task)
			continue
		}
		break
	}
	return tasks, nil
}

func (s *StoreMemory) GetNextRunTime() (time.Time, error) {
	if s.Tasks.Len() == 0 {
		return MaxDateTime, nil
	}
	for el := s.Tasks.Front(); el != nil; el = el.Next() {
		task := el.Value.(*Task)
		if task.Suspended {
			continue
		}
		return task.NextRunTime, nil
	}
	return MaxDateTime, nil
}

func (s *StoreMemory) AddTask(task *Task) error {
	if task.ID != 0 {
		return StoreInvalidTaskErr
	}
	atomic.AddInt64(&s.Index, 1)
	task.ID = s.Index
	for el := s.Tasks.Front(); el != nil; el = el.Next() {
		elTask := el.Value.(*Task)
		if task.NextRunTime.After(elTask.NextRunTime) {
			continue
		}
		s.TasksMap[task.ID] = s.Tasks.InsertBefore(task, el)
		return nil
	}
	s.TasksMap[task.ID] = s.Tasks.PushBack(task)
	return nil
}

func (s *StoreMemory) UpdateTask(task *Task) error {
	element, ok := s.TasksMap[task.ID]
	if !ok {
		return StoreNoTaskErr
	}
	element.Value = task
	for el := s.Tasks.Front(); el != nil; el = el.Next() {
		elTask := el.Value.(*Task)
		if elTask.ID == task.ID {
			continue
		}
		if task.NextRunTime.After(elTask.NextRunTime) {
			continue
		}
		s.Tasks.MoveBefore(element, el)
		return nil
	}
	s.Tasks.MoveToBack(element)
	return nil
}

func (s *StoreMemory) DelTask(task *Task) error {
	el, ok := s.TasksMap[task.ID]
	if !ok {
		return StoreNoTaskErr
	}
	delete(s.TasksMap, task.ID)
	s.Tasks.Remove(el)
	return nil
}

func (s *StoreMemory) GetTaskByID(id int64) (*Task, error) {
	el, ok := s.TasksMap[id]
	if !ok {
		return nil, StoreNoTaskErr
	}
	return el.Value.(*Task), nil
}
