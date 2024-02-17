package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"qsapi/internal"

	"github.com/gorilla/mux"
)

func main() {
	appContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := mux.NewRouter()

	db, err := internal.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	srv := internal.NewServer(appContext, router, db)

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
		cancel()
	case <-appContext.Done():
		log.Println("Application context canceled, exiting...")
	}
}
