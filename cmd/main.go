package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"qsapi/internal"
	"qsapi/pkg/pg_db"
	"qsapi/pkg/repo_cron"

	"github.com/gorilla/mux"
)

const (
	dsn = "user=user password=password dbname=quotation_db host=localhost port=5432 sslmode=disable"
)

func main() {
	appContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := mux.NewRouter()

	db, dbErr := pg_db.NewDatabasePg(dsn)
	panicIfErr(dbErr)

	ts := repo_cron.NewTaskScheduler()
	ts.Start()
	newTask := repo_cron.NewTask(db)
	cronErr := ts.AddTask("@every 30s", newTask.Run)
	panicIfErr(cronErr)

	srv := internal.NewServer(appContext, router, db, ts)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Listen and Serve failed: %v", err)
		}
	}()

	select {
	case <-stop:
		log.Println("Got SIGINT, exiting...")
		ts.Stop()
		cancel()
	case <-appContext.Done():
		log.Println("Application context canceled, exiting...")
	}
}
