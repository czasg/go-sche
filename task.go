package sche

import (
	"encoding/json"
	"github.com/czasg/gonal"
	"github.com/robfig/cron"
	"sync"
	"time"
)

var cronPool = sync.Map{}

type Trig string

func (t Trig) GetNextRunTime(previous time.Time) time.Time {
	cs, ok := cronPool.Load(t)
	if ok {
		return cs.(cron.Schedule).Next(previous)
	}
	ins, err := cron.Parse(string(t))
	if err != nil {
		return MaxDateTime
	}
	cronPool.Store(t, ins)
	return ins.Next(previous)
}

type Task struct {
	ID          int64             `json:"id" pg:",pk"`
	Name        string            `json:"name" pg:",use_zero"`
	Label       map[string]string `json:"label" pg:",use_zero"`
	Trig        Trig              `json:"trig" pg:",use_zero"`
	LastRunTime time.Time         `json:"last_run_time" pg:",use_zero"`
	NextRunTime time.Time         `json:"next_run_time" pg:",use_zero"`
	Suspended   bool              `json:"suspended" pg:",use_zero"`
}

func (t *Task) Run() error {
	t.LastRunTime = time.Now()
	t.NextRunTime = t.Trig.GetNextRunTime(t.LastRunTime)
	body, _ := json.Marshal(t)
	return gonal.Notify(nil, t.Label, body)
}
