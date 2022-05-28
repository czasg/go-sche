package sche

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func NewScheduler(stores ...Store) *Scheduler {
	var store Store
	if len(stores) > 0 {
		store = stores[0]
	}
	if store == nil {
		store = NewStoreMemory()
	}
	return &Scheduler{Store: store}
}

type Scheduler struct {
	Store  Store
	notify Notify
	lock   sync.Mutex
}

func (s *Scheduler) Start(ctx context.Context) error {
	if ctx == nil {
		return errors.New("ctx nil")
	}
	var wait time.Time
	s.notify = Notify{ctx: ctx}
	if s.Store == nil {
		s.Store = NewStoreMemory()
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.notify.Wait(wait):
		}
		now := time.Now()
		todo, err := s.Store.Todo(now)
		if err != nil {
			wait = MaxDateTime
			logrus.WithError(err).Error("获取待执行任务异常")
			continue
		}
		for _, t := range todo {
			log := logrus.WithField("taskName", t.Name)
			log.Info("触发任务")
			err = t.Run()
			if err != nil {
				log.WithError(err).Error("任务投放异常")
				continue
			}
			err = s.Store.UpdateTask(t)
			if err != nil {
				log.WithError(err).Error("更新任务状态异常")
			}
		}
		wait, err = s.Store.GetNextRunTime()
		if err != nil {
			wait = MaxDateTime
			logrus.WithError(err).Error("获取下一次调度时间异常")
		}
	}
}

func (s *Scheduler) AddTask(task *Task) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.Store == nil {
		s.Store = NewStoreMemory()
	}
	if task.Label == nil {
		task.Label = map[string]string{}
	}
	task.Label["task.label.name"] = task.Name
	task.NextRunTime = task.Trig.GetNextRunTime(time.Now())
	err := s.Store.AddTask(task)
	if err != nil {
		return err
	}
	s.notify.Notify()
	logrus.WithField("taskName", task.Name).Info("新增任务成功")
	return nil
}

func (s *Scheduler) UpdateTask(task *Task) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.Store == nil {
		s.Store = NewStoreMemory()
	}
	if task.Label == nil {
		task.Label = map[string]string{}
	}
	task.Label["task.label.name"] = task.Name
	task.NextRunTime = task.Trig.GetNextRunTime(time.Now())
	err := s.Store.UpdateTask(task)
	if err != nil {
		return err
	}
	s.notify.Notify()
	logrus.WithField("taskName", task.Name).Info("更新任务成功")
	return nil
}

func (s *Scheduler) DelTask(task *Task) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.Store == nil {
		s.Store = NewStoreMemory()
	}
	err := s.Store.DelTask(task)
	if err != nil {
		return err
	}
	s.notify.Notify()
	logrus.WithField("taskName", task.Name).Info("删除任务成功")
	return nil
}
