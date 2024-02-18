package repo_cron

import (
	"time"

	"github.com/robfig/cron/v3"
)

type ITaskScheduler interface {
	Start()
	Stop()
	AddTask(spec string, task func()) error
	GetResolveTime(id cron.EntryID) (time.Time, error)
	GetMainTaskID() cron.EntryID
}
