package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
	"github.com/MahikaJaguste/distributed-task-queue/pkg/worker"
)

func main() {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /execute", worker.HandleTaskExecution)

	db.SetupDb()

	port := 8001
	fmt.Printf("Server listening on %d!\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}

}
