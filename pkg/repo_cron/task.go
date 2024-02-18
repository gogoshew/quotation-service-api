package repo_cron

import (
	"qsapi/pkg/pg_db"

	"github.com/robfig/cron/v3"
)

type UpdateTask struct {
	ID cron.EntryID
	db *pg_db.DatabasePg
}

func NewTask(db *pg_db.DatabasePg) *UpdateTask {
	return &UpdateTask{db: db}
}

func (t *UpdateTask) Run() {
	t.db.GetRowsForUpdate()
}
