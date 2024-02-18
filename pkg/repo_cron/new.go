package repo_cron

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type TaskScheduler struct {
	taskID cron.EntryID
	cron   *cron.Cron
}

func NewTaskScheduler() ITaskScheduler {
	return &TaskScheduler{
		cron: cron.New(cron.WithLocation(time.UTC)),
	}
}

func (ts *TaskScheduler) AddTask(spec string, task func()) error {
	id, err := ts.cron.AddFunc(spec, task)
	if err != nil {
		log.Fatalf("Error add new task to cron: %v\n", err)
	}
	ts.taskID = id

	return nil
}

func (ts *TaskScheduler) GetResolveTime(id cron.EntryID) (time.Time, error) {
	resolveTime := ts.cron.Entry(id).Next
	if resolveTime.IsZero() {
		return time.Time{}, errors.New(fmt.Sprintf("Error there is not resolve time for task %v\n", id))
	}
	return resolveTime, nil
}

func (ts *TaskScheduler) Start() {
	ts.cron.Start()
}

func (ts *TaskScheduler) Stop() {
	ts.cron.Stop()
}

func (ts *TaskScheduler) GetMainTaskID() cron.EntryID {
	return ts.taskID
}
