package sche
//
//import (
//	"errors"
//	"github.com/go-pg/pg/v10"
//	"github.com/go-pg/pg/v10/orm"
//	"time"
//)
//
//var _ Store = (*PostgresStore)(nil)
//
//func NewPostgresStore(db *pg.DB) *PostgresStore {
//	err := db.Model((*Task)(nil)).CreateTable(&orm.CreateTableOptions{
//		IfNotExists: true,
//	})
//	if err != nil {
//		panic(err)
//	}
//	return &PostgresStore{db}
//}
//
//type PostgresStore struct {
//	DB *pg.DB
//}
//
//func (p *PostgresStore) Todo(now time.Time) ([]*Task, error) {
//	tasks := []*Task{}
//	err := p.DB.Model(&tasks).
//		Where("next_run_time <= ?", now).
//		Where("suspended = false").
//		Select()
//	if err != nil {
//		return nil, err
//	}
//	return tasks, nil
//}
//
//func (p *PostgresStore) GetNextRunTime() (time.Time, error) {
//	task := Task{}
//	err := p.DB.Model(&task).
//		Where("suspended = false").
//		Order("next_run_time ASC").
//		Returning("next_run_time").
//		Limit(1).
//		Select()
//	if err != nil {
//		if errors.Is(err, pg.ErrNoRows) {
//			return MaxDateTime, nil
//		}
//		return MaxDateTime, err
//	}
//	return task.NextRunTime, nil
//}
//
//func (p *PostgresStore) AddTask(task *Task) error {
//	if task.ID != 0 {
//		return StoreInvalidTaskErr
//	}
//	_, err := p.DB.Model(task).Insert()
//	return err
//}
//
//func (p *PostgresStore) UpdateTask(task *Task) error {
//	_, err := p.DB.Model(task).Where("id = ?id").Update()
//	return err
//}
//
//func (p *PostgresStore) DelTask(task *Task) error {
//	_, err := p.DB.Model(task).Where("id = ?id").Delete()
//	return err
//}
//
//func (p *PostgresStore) GetTaskByID(id int64) (*Task, error) {
//	task := Task{}
//	err := p.DB.Model(&task).Where("id = ?", id).Select()
//	if err != nil {
//		if errors.Is(err, pg.ErrNoRows) {
//			return nil, StoreNoTaskErr
//		}
//		return nil, err
//	}
//	return &task, nil
//}
