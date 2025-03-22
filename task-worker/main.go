package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	db "github.com/MahikaJaguste/distributed-task-queue/task-submission/database"
	te "github.com/MahikaJaguste/distributed-task-queue/task-worker/taskexecution"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /execute", te.HandleTaskExecution)

	db.SetupDb()

	port := 8001
	fmt.Printf("Server listening on %d!\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}

}
