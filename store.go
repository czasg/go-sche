package sche

import (
	"errors"
	"math"
	"time"
)

var (
	MaxDateTime = time.Now().Add(time.Duration(math.MaxInt64))
)

var (
	StoreInvalidTaskErr     = errors.New("store invalid task.")
	StoreNoTaskErr          = errors.New("store no task.")
)

type Store interface {
	Todo(now time.Time) ([]*Task, error)
	GetNextRunTime() (time.Time, error)
	AddTask(task *Task) error
	UpdateTask(task *Task) error
	DelTask(task *Task) error
	GetTaskByID(id int64) (*Task, error)
}
