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

func main() {
	appContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := mux.NewRouter()

	db, dbErr := pg_db.NewDatabasePg(os.Getenv("DB_CON_STR"))
	if dbErr != nil {
		panic(dbErr)
	}

	ts := repo_cron.NewTaskScheduler()
	ts.Start()
	newTask := repo_cron.NewTask(db)
	cronErr := ts.AddTask("@every 30s", newTask.Run)
	if dbErr != nil {
		panic(cronErr)
	}

	srv := internal.NewServer(appContext, router, db, ts)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Listen and Serve failed: %v", err)
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
