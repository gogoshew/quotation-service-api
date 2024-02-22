package repo_cron

import (
	"log"
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
	rows, err := t.db.GetRowsForUpdate()
	if err != nil {
		log.Printf("Error getting rows from cron task: %v\n", err)
		return
	}
	if len(rows) == 0 {
		log.Println("Nothing to update in cron task")
		return
	}

	for _, row := range rows {
		if err := t.db.UpdateLatest(row.Currency, row.Value); err != nil {
			log.Printf("Error updating currency value from cron task: %v\n", err)
		}
		log.Printf("Successfully updated value of %s\n", row.Currency)
	}

}
