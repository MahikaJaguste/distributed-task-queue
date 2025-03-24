package submission

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
	"github.com/joho/godotenv"

	forms "github.com/albrow/forms"
)

func StartSubmissionServer(port int) {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /submit", handleTaskSubmission)

	db.SetupDb()

	fmt.Printf("Submission server listening on %d!\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}

}

func handleTaskSubmission(w http.ResponseWriter, req *http.Request) {
	data, err := forms.Parse(req)
	if err != nil {
		// in case of any error
		// TODO
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := data.Get("name")

	result, err := db.DBCon.Exec("INSERT INTO tasks (name) VALUES (?)", name)
	if err != nil {
		// TODO
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		// TODO
		return
	}
	fmt.Printf("Task created with id = %d\n", id)
	w.Write(([]byte)(fmt.Sprintf("Task created with id = %d\n", id)))
}
