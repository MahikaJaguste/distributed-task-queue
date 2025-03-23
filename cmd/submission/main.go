package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
	"github.com/MahikaJaguste/distributed-task-queue/pkg/submission"
)

func main() {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /submit", submission.HandleTaskSubmission)

	db.SetupDb()

	port := 8080
	fmt.Printf("Server listening on %d!\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}

}
