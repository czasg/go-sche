package sche

import (
	"errors"
	"time"
)

var (
	MaxDateTime = time.Date(9999, 1, 1, 0, 0, 0, 0, time.Local)
)

var (
	StoreInvalidTaskErr = errors.New("store invalid task.")
	StoreNoTaskErr      = errors.New("store no task.")
)

type Store interface {
	Todo(now time.Time) ([]*Task, error)
	GetNextRunTime() (time.Time, error)
	AddTask(task *Task) error
	UpdateTask(task *Task) error
	DelTask(task *Task) error
	GetTaskByID(id int64) (*Task, error)
}
